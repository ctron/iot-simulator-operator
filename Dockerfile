FROM centos:7

MAINTAINER Jens Reimann <jreimann@redhat.com>
LABEL maintainer="Jens Reimann <jreimann@redhat.com>"

RUN yum -y update
RUN yum -y install epel-release

ENV \
    GOPATH=/go

RUN mkdir -p /go/src/github.com/ctron
ADD . /go/src/github.com/ctron/iot-simulator-operator

RUN yum -y install golang && \
    go version && \
    cd /go/src/github.com/ctron/iot-simulator-operator/cmd/manager && go build -o /iot-simulator-operator . && \
    cd / && \
    rm -Rf go && \
    yum -y history undo last && yum -y clean all && \
    true

ENV IMAGE_TAG=":latest"

ENTRYPOINT /iot-simulator-operator
