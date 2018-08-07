# go-splunk-ec-client

Splunk's HTTP Event Collector (EC) client. The EC is an endpoint allowing sending
messages to Splunk via HTTP. The endpoint identifies its clients based on a token
the clients' provide. A Splunk administrator configures tokens under "Add Data",
"HTTP Event Collector". Once configured, the administrator provides the token
to a client application, e.g. `go-splunk-ec-client`.

The "Input Settings" for the collector are:
* Source Type: Automatic
* App context: Search & Reporting
* Index: `main`

By default, the HTTP Event Collector receives data over HTTPS on TCP port 8088. 

## Getting Started

```golang
package main

import (
    "fmt"
    "time"
    splunk "github.com/greenpau/go-splunk-ec-client"
)

func main() {
    host := "localhost"
    port := 8088
    token := "61876693-4758-4f45-bca7-c910ccc746eb"
    cli, err := splunk.NewClient(host, port, token)
    if err != nil {
        return err
    }
    msg := fmt.Sprintf("{\"message\": "test message on %s\"}", time.Now().String())
    if err := cli.Send(msg); err != nil {
        return err
    }
}
```
