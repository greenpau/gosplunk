package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Event represents an event sent to HTTP Event Collector. It conforms to the
// standard described in `services/collector` REST API documentation at
// http://docs.splunk.com/Documentation/Splunk/7.1.2/RESTREF/RESTinput#services.2Fcollector
type Event struct {
	Channel    string            `json:"channel,omitempty" yaml:"channel"`
	Message    string            `json:"event" yaml:"event"`
	Fields     map[string]string `json:"fields,omitempty" yaml:"fields"`
	Host       string            `json:"host,omitempty" yaml:"host"`
	Index      string            `json:"index,omitempty" yaml:"index"`
	Source     string            `json:"source,omitempty" yaml:"source"`
	SourceType string            `json:"sourcetype,omitempty" yaml:"sourcetype"`
	Time       uint64            `json:"time,omitempty" yaml:"time"`
}

// EventResponse is the response payload of an event submission to HEC
// endpoint. The response ordinarily contains the following fields:
//
//   * Text: Human readable status, same value as code.
//   * Code: Machine format status, same value as text.
//   * InvalidEvent: This field gets populated when errors occur. It indicates
//     the zero-based index of first invalid event in an event sequence.
//   * AckID: This field gets populated when "useACK" is enabled for a token.
//     It indicates the "ackId" to use for checking an indexer acknowledgement.
//
// The following helps understanding the meaning of the values of the Code field:
//   * 0: "200 OK", Success
//   * 1: "403 Forbidden", Token disabled
//   * 2: "401 Unauthorized", Token is required
//   * 3: "401 Unauthorized", Invalid authorization
//   * 4: "403 Forbidden", Invalid token
//   * 5: "400 Bad Request", No data
//   * 6: "400 Bad Request", Invalid data format
//   * 7: "400 Bad Request", Incorrect index
//   * 8: "500 Internal Error", Internal server error
//   * 9: "503 Service Unavailable", Server is busy
//   * 10: "400 Bad Request", Data channel is missing
//   * 11: "400 Bad Request", Invalid data channel
//   * 12: "400 Bad Request", Event field is required
//   * 13: "400 Bad Request", Event field cannot be blank
//   * 14: "400 Bad Request", ACK is disabled
//   * 15: "400 Bad Request", Error in handling indexed fields
//   * 16: "400 Bad Request", Query string authorization is not enabled
type EventResponse struct {
	Text         string `json:"text" yaml:"text"`
	Code         int    `json:"code" yaml:"code"`
	InvalidEvent int    `json:"invalid-event-number" yaml:"invalid-event-number"`
	AckID        int    `json:"ackId" yaml:"ackId"`
}

// Send function sends events to HTTP Event Collector using the Splunk platform
// JSON event protocol. The interface is described at
// http://docs.splunk.com/Documentation/Splunk/7.1.2/RESTREF/RESTinput#services.2Fcollector
func (cli *Client) Send(evt Event) error {
	log.Debugf("%s: url=%s", cli.Name, cli.Endpoints.Event)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(evt)
	req, err := http.NewRequest("POST", cli.Endpoints.Event, b)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", cli.Token))
	resp, err := cli.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: send() failed: %s", cli.Name, err)
	}
	defer resp.Body.Close()
	log.Debugf("%s: status=%d %s", cli.Name, resp.StatusCode, http.StatusText(resp.StatusCode))

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: send() failed: erred while reading response body: %s", cli.Name, err)
	}
	eventResponse := new(EventResponse)
	if err := json.Unmarshal(respBody, eventResponse); err != nil {
		return fmt.Errorf("%s: send() failed: the response is not JSON but: %s", cli.Name, respBody)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: send() failed: %d %s (%s)", cli.Name, resp.StatusCode, http.StatusText(resp.StatusCode), eventResponse.Text)
	}
	log.Debugf("%s: code=%d, text=%s", cli.Name, eventResponse.Code, eventResponse.Text)
	return nil
}
