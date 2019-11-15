package main

// configuration module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	logs "github.com/sirupsen/logrus"
)

// Configuration stores server configuration parameters
type Configuration struct {
	Port         int    `json:"port"`         // server port number
	ExitCodes    string `json:"exitCodes"`    // location of exit codes
	ExitCodesUrl string `json:"exitCodesUrl"` // URL location of exit codes
	CacheExpire  int64  `json:"cacheExpire"`  // cache expiration in sec
	Templates    string `json:"templates"`    // location of server templates
	Jscripts     string `json:"jscripts"`     // location of server JavaScript files
	Images       string `json:"images"`       // location of server images
	Styles       string `json:"styles"`       // location of server CSS styles
	LogFormatter string `json:"logFormatter"` // LogFormatter type
	Verbose      int    `json:"verbose"`      // verbosity level
	ServerKey    string `json:"ckey"`         // tls.key file
	ServerCrt    string `json:"cert"`         // tls.cert file
	RootCA       string `json:"rootCA"`       // RootCA file
}

// Config variable represents configuration object
var Config Configuration

// String returns string representation of server Config
func (c *Configuration) String() string {
	data, _ := json.Marshal(c)
	return fmt.Sprintf(string(data))
}

// ParseConfig parse given config file
func ParseConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		logs.WithFields(logs.Fields{"configFile": configFile}).Fatal("Unable to read", err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		logs.WithFields(logs.Fields{"configFile": configFile}).Fatal("Unable to parse", err)
		return err
	}
	return nil
}
