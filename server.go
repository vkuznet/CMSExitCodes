package main

// web server module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	logs "github.com/sirupsen/logrus"

	_ "expvar"         // to be used for monitoring, see https://github.com/divan/expvarmon
	_ "net/http/pprof" // profiler, see https://golang.org/pkg/net/http/pprof/
)

// global variables
var _top, _bottom, _search string
var _cache map[string]Cache

// Time0 represents initial time when we started the server
var Time0 time.Time

// Server code
func Server(configFile string) {
	Time0 = time.Now()
	err := ParseConfig(configFile)
	if Config.LogFormatter == "json" {
		logs.SetFormatter(&logs.JSONFormatter{})
	} else if Config.LogFormatter == "text" {
		logs.SetFormatter(&logs.TextFormatter{})
	} else {
		logs.SetFormatter(&logs.JSONFormatter{})
	}
	if Config.Verbose > 0 {
		logs.SetLevel(logs.DebugLevel)
	}
	if err != nil {
		logs.WithFields(logs.Fields{"Time": time.Now(), "Config": configFile}).Error("Unable to parse")
	}
	fmt.Println("Configuration:", Config.String())

	// initialize cache variable
	_cache = make(map[string]Cache)

	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	tmplData["Time"] = time.Now()
	tmplData["Version"] = info()
	tmplData["Base"] = Config.Base
	_top = templates.Top(Config.Templates, tmplData)
	_bottom = templates.Bottom(Config.Templates, tmplData)

	// assign handlers
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(Config.Styles))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(Config.Jscripts))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(Config.Images))))
	http.HandleFunc("/", AuthHandler)

	// Start server
	addr := fmt.Sprintf(":%d", Config.Port)
	_, e1 := os.Stat(Config.ServerCrt)
	_, e2 := os.Stat(Config.ServerKey)
	if e1 == nil && e2 == nil {
		//start HTTPS server which require user certificates
		rootCA := x509.NewCertPool()
		caCert, _ := ioutil.ReadFile(Config.RootCA)
		rootCA.AppendCertsFromPEM(caCert)
		server := &http.Server{
			Addr: addr,
			TLSConfig: &tls.Config{
				//                 ClientAuth: tls.RequestClientCert,
				RootCAs: rootCA,
			},
		}
		logs.WithFields(logs.Fields{"Addr": addr}).Info("Starting HTTPs server")
		err = server.ListenAndServeTLS(Config.ServerCrt, Config.ServerKey)
	} else {
		// Start server without user certificates
		logs.WithFields(logs.Fields{"Addr": addr}).Info("Starting HTTP server")
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		logs.WithFields(logs.Fields{
			"Error": err,
		}).Fatal("ListenAndServe: ")
	}
}
