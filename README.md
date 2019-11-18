### CMS Exit Codes service

[![Build Status](https://travis-ci.org/vkuznet/CMSExitCodes.svg?branch=master)](https://travis-ci.org/vkuznet/CMSExitCodes)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkuznet/CMSExitCodes)](https://goreportcard.com/report/github.com/vkuznet/CMSExitCodes)
[![GoDoc](https://godoc.org/github.com/vkuznet/CMSExitCodes?status.svg)](https://godoc.org/github.com/vkuznet/CMSExitCodes)

This directory contains codebase for CMS Exit Codes service.
The service returns all known exit codes from different CMS sub-systems, e.g.
WMAgent, JobExit, StageOut, etc. All codes are defined in `codes` area.
The service can either return full list or a specific code, e.g.
here is a client side:
```
curl http://xxx.yyy.com/exitcodes/8021
curl -H "Accept: application/json" http://xxx.yyy.com/exitcodes/8021
{"JobExit":{"8021":"FileReadError (May be a site error)"},"WMAgent":{"8021":"FileReadError (May be a site error)"}}
```

### Setup

To build it please install Go language on your system
and series of dependencies:

```
# obtain necessary dependencies
go get github.com/sirupsen/logrus
go get github.com/shirou/gopsutil
go get github.com/divan/expvarmon
go get github.com/sirupsen/logrus

# build server
go build
```

To run the service use the following command:
```
CMSExitCodes -config server.json
```
where `server.json` has the following form:
```
{
    "exitCodes":"/path/exitCodes.txt",
    "exitCodesUrl": "https://raw.githubusercontent.com/vkuznet/CMSExitCodes/master/codes",
    "cacheExpire": 600,
    "port": 9201,
    "templates": "/path/templates",
    "jscripts": "/path/js",
    "styles": "/path/css",
    "images": "/path/images",
    "verbose": 0
}
```
