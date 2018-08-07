package client

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

// Configuration is a configuration for the Client. It allows specifying
// host and port of a Splunk HEC endpoint, as well as HEC token.
type Configuration struct {
	LogLevel  log.Level `json:"log_level" yaml:"log_level"`
	File      string    `json:"conf_file" yaml:"conf_file"`
	Collector struct {
		Proto   string `json:"proto" yaml:"proto"`
		Host    string `json:"host" yaml:"host"`
		Port    int    `json:"port" yaml:"port"`
		Token   string `json:"token" yaml:"token"`
		Timeout int
	}
}

// NewConfiguration creates an instance of a Configuration from the
// configuration file in YAML format on a local filesystem.
func NewConfiguration(f string) (Configuration, error) {
	c := Configuration{}
	filePath := ""
	if strings.HasPrefix(f, "~/") {
		usr, err := user.Current()
		if err != nil {
			return c, err
		}
		f = filepath.Join(usr.HomeDir, f[2:])
	}
	filePath, err := filepath.Abs(f)
	if err != nil {
		return c, err
	}
	if c.File == "" {
		c.File = filePath
	}
	content, err := ioutil.ReadFile(c.File)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return c, err
	}
	if c.Collector.Timeout == 0 {
		c.Collector.Timeout = 5
	}
	if c.LogLevel == 0 {
		// Default: INFO
		log.SetLevel(4)
		c.LogLevel = log.GetLevel()
	} else {
		log.SetLevel(c.LogLevel)
	}
	if c.Collector.Host == "" {
		return c, errors.New("collector.host is undefined")
	}
	if c.Collector.Token == "" {
		return c, errors.New("collector.token is undefined")
	}
	if c.Collector.Port == 0 {
		c.Collector.Port = 8088
	}
	if c.Collector.Proto == "" {
		c.Collector.Proto = "https"
	}
	if c.Collector.Proto != "http" && c.Collector.Proto != "https" {
		return c, errors.New("collector.proto must be either http or https")
	}
	return c, nil
}
