package client

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	testFailed := 0

	healthCheckHandler := func(w http.ResponseWriter, r *http.Request) {
		resp := HealthCheckResponse{200, "Success"}
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		t.Fatalf("Failed to start a web server: %s", err)
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(healthCheckHandler))

	server.Listener.Close()
	server.Listener = listener
	server.Start()

	defer server.Close()

	log.Infof("%s", server.URL)

	for i, test := range []struct {
		configFile string
		shouldFail bool
	}{
		{configFile: "../../examples/.splunk.hec.yaml", shouldFail: false},
	} {
		conf, err := NewConfiguration(test.configFile)
		if err != nil {
			t.Logf("FAIL: Test %d: configuration file '%s', expected to pass, but failed with: %v", i, test.configFile, err)
			testFailed++
			continue
		}
		conf.LogLevel = log.DebugLevel
		log.SetLevel(conf.LogLevel)
		cli, err := NewClient(conf)
		if err != nil {
			if !test.shouldFail {
				t.Logf("FAIL: Test %d: configuration file '%s', expected to pass, but failed with: %v", i, test.configFile, err)
				testFailed++
				continue
			}
			t.Logf("PASS: Test %d: configuration file '%s', expected to fail: failed with: %v", i, test.configFile, err)
			continue
		}
		if test.shouldFail {
			t.Logf("FAIL: Test %d: configuration file '%s', expected to fail, but passed", i, test.configFile)
			testFailed++
			continue
		}
		t.Logf("PASS: Test %d: configuration file '%s', expected to pass: passed, token: %s", i, test.configFile, cli.Token)
	}
	if testFailed > 0 {
		t.Fatalf("Failed %d tests", testFailed)
	}
}
