package main

// templates module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"bytes"
	"html/template"
	"path/filepath"
)

// consume list of templates and release their full path counterparts
func fileNames(tdir string, filenames ...string) []string {
	flist := []string{}
	for _, fname := range filenames {
		flist = append(flist, filepath.Join(tdir, fname))
	}
	return flist
}

// parse template with given data
func parseTmpl(tdir, tmpl string, data interface{}) string {
	buf := new(bytes.Buffer)
	filenames := fileNames(tdir, tmpl)
	funcMap := template.FuncMap{
		// The name "oddFunc" is what the function will be called in the template text.
		"oddFunc": func(i int) bool {
			if i%2 == 0 {
				return true
			}
			return false
		},
		// The name "inListFunc" is what the function will be called in the template text.
		"inListFunc": func(a string, list []string) bool {
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
		},
	}
	t := template.Must(template.New(tmpl).Funcs(funcMap).ParseFiles(filenames...))
	err := t.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// ServerTemplates structure
type ServerTemplates struct {
	top, bottom, faq, codes, status string
}

// Top method for ServerTemplates structure
func (q ServerTemplates) Top(tdir string, tmplData map[string]interface{}) string {
	if q.top != "" {
		return q.top
	}
	q.top = parseTmpl(Config.Templates, "top.tmpl", tmplData)
	return q.top
}

// Bottom method for ServerTemplates structure
func (q ServerTemplates) Bottom(tdir string, tmplData map[string]interface{}) string {
	if q.bottom != "" {
		return q.bottom
	}
	q.bottom = parseTmpl(Config.Templates, "bottom.tmpl", tmplData)
	return q.bottom
}

// FAQ method for ServerTemplates structure
func (q ServerTemplates) FAQ(tdir string, tmplData map[string]interface{}) string {
	if q.faq != "" {
		return q.faq
	}
	q.faq = parseTmpl(Config.Templates, "faq.tmpl", tmplData)
	return q.faq
}

// Codes method for ServerTemplates structure
func (q ServerTemplates) Codes(tdir string, tmplData map[string]interface{}) string {
	if q.codes != "" {
		return q.codes
	}
	q.codes = parseTmpl(Config.Templates, "codes.tmpl", tmplData)
	return q.codes
}

// Status method for ServerTemplates structure
func (q ServerTemplates) Status(tdir string, tmplData map[string]interface{}) string {
	if q.status != "" {
		return q.status
	}
	q.status = parseTmpl(Config.Templates, "status.tmpl", tmplData)
	return q.status
}
