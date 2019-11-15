package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// helper functions for handlers module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//
type Cache struct {
	Name       string
	Codes      map[string]interface{}
	LastUpdate int64
}

// helper function to read exit codes from HTTP URL location
func exitCodesHttpReader(url string) (map[string]interface{}, error) {
	codes := make(map[string]interface{})
	// read exit code file
	resp, err := http.Get(url)
	if err != nil {
		return codes, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		arr := strings.Split(line, ":")
		codes[strings.Trim(arr[0], " ")] = strings.Trim(arr[1], " ")
	}
	return codes, nil
}

// helper function to read exit codes file
func exitCodesReader(fname string) (map[string]interface{}, error) {
	codes := make(map[string]interface{})
	// read exit code file
	fin := fmt.Sprintf("%s/%s", Config.ExitCodes, fname)
	data, err := ioutil.ReadFile(fin)
	if err != nil {
		return codes, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		arr := strings.Split(line, ":")
		codes[strings.Trim(arr[0], " ")] = strings.Trim(arr[1], " ")
	}
	return codes, nil
}

func exitCodesCache(name string) (map[string]interface{}, error) {
	// check cache first
	if cache, ok := _cache[name]; ok {
		if cache.Name == name && cache.Codes != nil && time.Now().Unix()-cache.LastUpdate < Config.CacheExpire {
			log.Println("read from cache", name)
			return cache.Codes, nil
		}
	}
	var codes map[string]interface{}
	var err error
	if Config.ExitCodesUrl != "" {
		arr := strings.Split(name, "/")
		fname := arr[len(arr)-1]
		url := fmt.Sprintf("%s/%s", Config.ExitCodesUrl, fname)
		log.Println("read from Http", url)
		codes, err = exitCodesHttpReader(url)
	} else {
		log.Println("read from file", name)
		codes, err = exitCodesReader(name)
	}
	// put codes into cache and then return them back
	_cache[name] = Cache{Name: name, Codes: codes, LastUpdate: time.Now().Unix()}
	return codes, err
}

// helper function to read exit codes file
func exitCodes(fname string) (map[string]interface{}, error) {
	codes := make(map[string]interface{})
	// read exit code file
	if fname == "" || strings.ToLower(fname) == "all" {
		file, err := os.Open(Config.ExitCodes)
		defer file.Close()
		if err != nil {
			return codes, err
		}
		list, _ := file.Readdirnames(0) // read all files/dir
		for _, name := range list {
			c, e := exitCodesCache(name)
			if e != nil {
				return codes, e
			}
			system := strings.Split(name, ".txt")[0]
			codes[system] = c
		}
		return codes, nil
	}
	file, err := os.Open(Config.ExitCodes)
	defer file.Close()
	if err != nil {
		return codes, err
	}
	list, _ := file.Readdirnames(0) // read all files/dir
	for _, name := range list {
		if strings.Contains(strings.ToLower(name), strings.ToLower(fname)) {
			return exitCodesCache(name)
		}
	}
	return codes, nil
}

// helper function to extract username from auth-session cookie
func username(r *http.Request) (string, error) {
	user := "test"
	return user, nil
}

// authentication function
func auth(r *http.Request) error {
	_, err := username(r)
	return err
}
