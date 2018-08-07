package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// HealthCheckResponse is the response payload of a health check. It checks
// whether there is space available in the queue. Per the specification, the
// value of Code field of HealthCheckResponse has the following meaning:
//
//   * 200: HEC is available and accepting input
//   * 400: Invalid HEC token
//   * 503: HEC is unhealthy, queues are full
type HealthCheckResponse struct {
	Code int    `json:"code" yaml:"code"`
	Text string `json:"text" yaml:"text"`
}

// HealthCheck performs a health check according to Splunk's REST API
// specification at http://docs.splunk.com/Documentation/Splunk/7.1.2/RESTREF/RESTinput#services.2Fcollector.2Fhealth.
func (cli *Client) HealthCheck() error {
	log.Debugf("%s: url=%s", cli.Name, cli.Endpoints.Health)
	req, err := http.NewRequest("GET", cli.Endpoints.Health, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", cli.Token))
	resp, err := cli.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: health check failed: %s", cli.Name, err)
	}
	defer resp.Body.Close()
	log.Debugf("%s: status=%d %s", cli.Name, resp.StatusCode, http.StatusText(resp.StatusCode))
	switch resp.StatusCode {
	case 200:
		log.Debugf("%s: HEC is available and accepting input", cli.Name)
	case 400:
		return fmt.Errorf("%s: health check failed: Invalid HEC token", cli.Name)
	case 503:
		return fmt.Errorf("%s: health check failed: HEC is unhealthy, queues are full", cli.Name)
	default:
		return fmt.Errorf("%s: health check failed: %d %s", cli.Name, resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: health check failed: erred while reading response body: %s", cli.Name, err)
	}
	healthCheckResponse := new(HealthCheckResponse)
	if err := json.Unmarshal(respBody, healthCheckResponse); err != nil {
		return fmt.Errorf("%s: health check failed: the response is not JSON but: %s", cli.Name, respBody)
	}
	log.Debugf("%s: code=%d, text=%s", cli.Name, healthCheckResponse.Code, healthCheckResponse.Text)
	return nil
}
