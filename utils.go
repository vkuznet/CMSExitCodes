package main

// utils module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

import (
	"fmt"
	"sort"

	logs "github.com/sirupsen/logrus"
)

// ErrPropagate error helper function which can be used in defer ErrPropagate()
func ErrPropagate(api string) {
	if err := recover(); err != nil {
		logs.WithFields(logs.Fields{
			"api":   api,
			"error": err,
		}).Error("Server ERROR")
		panic(fmt.Sprintf("%s:%s", api, err))
	}
}

// FindInList helper function to find item in a list
func FindInList(a string, arr []string) bool {
	for _, e := range arr {
		if e == a {
			return true
		}
	}
	return false
}

// InList helper function to check item in a list
func InList(a string, list []string) bool {
	check := 0
	for _, b := range list {
		if b == a {
			check += 1
		}
	}
	if check != 0 {
		return true
	}
	return false
}

// MapKeys helper function to return keys from a map
func MapKeys(rec map[int]string) []int {
	keys := make([]int, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	sort.Sort(IntList(keys))
	return keys
}

// IntList implement sort for []int type
type IntList []int

// Len provides length of the []int type
func (s IntList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s IntList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s IntList) Less(i, j int) bool { return s[i] < s[j] }

// StringList implement sort for []string type
type StringList []string

// Len provides length of the []int type
func (s StringList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s StringList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }
