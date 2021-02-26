FROM golang:1.15

WORKDIR /go/src

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLE=1

RUN apt-get update && \
	apt-get install build-essential -y 

cmd ["tail", "-f", "/dev/null"]
