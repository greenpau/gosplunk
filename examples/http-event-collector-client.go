package main

import (
	"fmt"
	splunk "github.com/greenpau/gosplunk/http-event-collector/client"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	configFile := "~/.splunk.hec.yaml"
	conf, err := splunk.NewConfiguration(configFile)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	conf.LogLevel = log.DebugLevel
	log.SetLevel(conf.LogLevel)
	cli, err := splunk.NewClient(conf)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	msg := splunk.Event{
		Message: fmt.Sprintf("test message on %s\"", time.Now().String()),
		Fields: map[string]string{
			"foo": "bar",
			"bar": "foo",
		},
	}
	log.Debugf("message=\"%v\"", msg)
	if err := cli.Send(msg); err != nil {
		log.Errorf("%s", err)
		return
	}
}
