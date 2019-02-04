FROM fedora:29

MAINTAINER Jens Reimann <jreimann@redhat.com>
LABEL maintainer="Jens Reimann <jreimann@redhat.com>"

RUN dnf -y update
RUN dnf -y install golang
RUN go version

ENV \
    GOPATH=/go

RUN mkdir -p /go/src/github.com/ctron
ADD . /go/src/github.com/ctron/iot-simulator-operator

WORKDIR /go/src/github.com/ctron/iot-simulator-operator

RUN cd cmd/manager && go build -o /iot-simulator-operator .

WORKDIR /

ENTRYPOINT /iot-simulator-operator
