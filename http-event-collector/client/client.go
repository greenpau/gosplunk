// Package client implements Splunk's HTTP Event Collector (EC) client.
package client

import (
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

// Client is a Splunk HEC client. It sends messages to Splunk's
// RESTful API using HTTP/S transport.
//
// The Client uses HEC Tokens to authenticate to the API.
type Client struct {
	mux       sync.Mutex
	client    *http.Client
	Name      string
	Endpoints struct {
		Health string
		Event  string
		Raw    string
	}
	Token string
}

// NewClient initiates an instance of a Client based on the
// Configuration provided. The function populates URL for various
// endpoints, e.g. health, event, etc. Lastly, the function performs
// a health check to assess whether HEC interfaces is available.
// Upon successful completion of the health chech, the function
// returns an instance of the Client.
func NewClient(c Configuration) (Client, error) {
	cli := Client{
		Name: "splunk-http-collector-client",
	}
	if err := cli.Configure(c.Collector.Proto, c.Collector.Host, c.Collector.Port); err != nil {
		return cli, err
	}
	log.Debugf("%s: proto=%s", cli.Name, c.Collector.Proto)
	log.Debugf("%s: host=%s", cli.Name, c.Collector.Host)
	log.Debugf("%s: port=%d", cli.Name, c.Collector.Port)
	log.Debugf("%s: token=%s", cli.Name, c.Collector.Token)
	log.Debugf("%s: timeout=%d", cli.Name, c.Collector.Timeout)
	log.Debugf("%s: endpoint.health=%s", cli.Name, cli.Endpoints.Health)
	log.Debugf("%s: endpoint.event=%s", cli.Name, cli.Endpoints.Event)
	log.Debugf("%s: endpoint.raw=%s", cli.Name, cli.Endpoints.Raw)
	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	cli.client = &http.Client{
		Timeout:   time.Duration(c.Collector.Timeout) * time.Second,
		Transport: t,
	}
	cli.Token = c.Collector.Token
	if err := cli.HealthCheck(); err != nil {
		return cli, err
	}
	return cli, nil
}

// Configure function populates URL for various endpoints, e.g. health,
// event, etc.
func (cli *Client) Configure(proto string, host string, port int) error {
	cli.Endpoints.Health = fmt.Sprintf("%s://%s:%d/services/collector/health", proto, host, port)
	cli.Endpoints.Event = fmt.Sprintf("%s://%s:%d/services/collector/event", proto, host, port)
	cli.Endpoints.Raw = fmt.Sprintf("%s://%s:%d/services/collector/raw", proto, host, port)
	return nil
}
