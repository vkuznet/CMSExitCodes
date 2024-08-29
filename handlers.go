package main

// handlers module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	logs "github.com/sirupsen/logrus"
)

// TotalGetRequests counts total number of GET requests received by the server
var TotalGetRequests uint64

// TotalPostRequests counts total number of POST requests received by the server
var TotalPostRequests uint64

// ServerSettings controls server parameters
type ServerSettings struct {
	Level        int    `json:"level"`        // verbosity level
	LogFormatter string `json:"logFormatter"` // logrus formatter
}

//
// ### HTTP METHODS
//

// AuthHandler authenticate incoming requests and route them to appropriate handler
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// increment GET/POST counters
	if r.Method == "GET" {
		atomic.AddUint64(&TotalGetRequests, 1)
	}
	if r.Method == "POST" {
		atomic.AddUint64(&TotalPostRequests, 1)
	}

	// check if server started with hkey file (auth is required)
	var user string
	var err error
	user, err = username(r)
	if err != nil {
		return
	}
	logs.WithFields(logs.Fields{"User": user, "Path": r.URL.Path}).Info("")
	// define all methods which requires authentication
	arr := strings.Split(r.URL.Path, "/")
	path := arr[len(arr)-1]
	switch path {
	case "faq":
		FAQHandler(w, r)
	case "status":
		StatusHandler(w, r)
	case "server":
		SettingsHandler(w, r)
	default:
		DataHandler(w, r)
	}
}

// GET Methods

// DataHandler handlers Data requests
func DataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// check HTTP header
	var accept, content string
	if _, ok := r.Header["Accept"]; ok {
		accept = r.Header["Accept"][0]
	}
	if _, ok := r.Header["Content-Type"]; ok {
		content = r.Header["Content-Type"][0]
	}
	user, _ := username(r)
	name := r.FormValue("name")
	codes, err := exitCodes(name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("unable to get exit codes, error=%v", err)))
		return
	}
	if Config.Verbose > 0 {
		log.Println("user", user, "name", name, "found", len(codes), "codes files")
	}
	arr := strings.Split(r.URL.Path, "/")
	key := arr[len(arr)-1]
	if len(arr) > 1 {
		cfinal := make(map[string]interface{})
		for k, val := range codes {
			switch vvv := val.(type) {
			case map[string]interface{}:
				if key != "" {
					if v, ok := vvv[key]; ok {
						c := make(map[string]string)
						c[key] = v.(string)
						cfinal[k] = c
					}
				} else {
					cfinal[k] = val
				}
			case string:
				if k == key {
					cfinal[k] = val
				}
			}
		}
		codes = cfinal
	}
	// return json representation of the data
	if strings.Contains(accept, "json") || strings.Contains(content, "json") {
		data, err := json.Marshal(codes)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("unable to marshal data, error=%v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	// return web UI representation of the data
	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	tmplData["User"] = user
	tmplData["Date"] = time.Now().Unix()
	tmplData["Name"] = name
	tmplData["ExitCodes"] = codes
	page := templates.Codes(Config.Templates, tmplData)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(_top + page + _bottom))
}

// FAQHandler handlers FAQ requests
func FAQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	page := templates.FAQ(Config.Templates, tmplData)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(_top + page + _bottom))
}

// Memory structure keeps track of server memory
type Memory struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// Mem structure keeps track of virtual/swap memory of the server
type Mem struct {
	Virtual Memory
	Swap    Memory
}

// StatusHandler handlers Status requests
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// check HTTP header
	var accept, content string
	if _, ok := r.Header["Accept"]; ok {
		accept = r.Header["Accept"][0]
	}
	if _, ok := r.Header["Content-Type"]; ok {
		content = r.Header["Content-Type"][0]
	}

	// get cpu and mem profiles
	m, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	l, _ := load.Avg()
	c, _ := cpu.Percent(time.Millisecond, true)
	process, perr := process.NewProcess(int32(os.Getpid()))

	// get unfinished queries
	var templates ServerTemplates
	tmplData := make(map[string]interface{})
	tmplData["NGo"] = runtime.NumGoroutine()
	//     virt := Memory{Total: m.Total, Free: m.Free, Used: m.Used, UsedPercent: m.UsedPercent}
	//     swap := Memory{Total: s.Total, Free: s.Free, Used: s.Used, UsedPercent: s.UsedPercent}
	tmplData["Memory"] = m.UsedPercent
	tmplData["Swap"] = s.UsedPercent
	tmplData["Load1"] = l.Load1
	tmplData["Load5"] = l.Load5
	tmplData["Load15"] = l.Load15
	tmplData["CPU"] = c
	if perr == nil { // if we got process info
		conn, err := process.Connections()
		if err == nil {
			tmplData["Connections"] = conn
		}
		openFiles, err := process.OpenFiles()
		if err == nil {
			tmplData["OpenFiles"] = openFiles
		}
	}
	tmplData["Uptime"] = time.Since(Time0).Seconds()
	tmplData["GetRequests"] = TotalGetRequests
	tmplData["PostRequests"] = TotalPostRequests
	page := templates.Status(Config.Templates, tmplData)
	if strings.Contains(accept, "json") || strings.Contains(content, "json") {
		data, err := json.Marshal(tmplData)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("unable to marshal data, error=%v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(_top + page + _bottom))
}

// POST methods

// SettingsHandler handlers Settings requests
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var s = ServerSettings{}
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if s.LogFormatter == "json" {
		logs.SetFormatter(&logs.JSONFormatter{})
	} else if s.LogFormatter == "text" {
		logs.SetFormatter(&logs.TextFormatter{})
	} else {
		logs.SetFormatter(&logs.TextFormatter{})
	}
	logs.WithFields(logs.Fields{
		"Verbose level": s.Level,
		"Log formatter": s.LogFormatter,
	}).Debug("update server settings")
	w.WriteHeader(http.StatusOK)
	return
}
