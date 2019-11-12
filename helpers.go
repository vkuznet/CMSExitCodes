package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// helper functions for handlers module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//

// helper function to read exit codes file
func exitCodesReader(fname string) (map[int]string, error) {
	codes := make(map[int]string)
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
		code, _ := strconv.Atoi(arr[0])
		codes[code] = strings.Trim(arr[1], " ")
	}
	return codes, nil
}

// helper function to read exit codes file
func exitCodes(fname string) (map[int]string, error) {
	codes := make(map[int]string)
	// read exit code file
	if fname == "" || strings.ToLower(fname) == "all" {
		file, err := os.Open(Config.ExitCodes)
		defer file.Close()
		if err != nil {
			return codes, err
		}
		list, _ := file.Readdirnames(0) // read all files/dir
		for _, name := range list {
			c, e := exitCodesReader(name)
			if e != nil {
				return codes, e
			}
			for k, v := range c {
				codes[k] = v
			}
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
			return exitCodesReader(name)
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