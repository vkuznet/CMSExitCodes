# Start from the latest golang base image
FROM golang:latest
MAINTAINER Valentin Kuznetsov vkuznet@gmail.com

ENV WDIR=/data
WORKDIR $WDIR

RUN go get github.com/sirupsen/logrus
RUN go get -d github.com/shirou/gopsutil/...
RUN go get github.com/vkuznet/CMSExitCodes
RUN sed -i -e "s,ROOT,/go/src/github.com/vkuznet/CMSExitCodes,g" /go/src/github.com/vkuznet/CMSExitCodes/server.json
CMD ["CMSExitCodes", "-config", "/go/src/github.com/vkuznet/CMSExitCodes/server.json"]
