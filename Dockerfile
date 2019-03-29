FROM fedora:29

MAINTAINER Jens Reimann <jreimann@redhat.com>
LABEL maintainer="Jens Reimann <jreimann@redhat.com>"

RUN dnf -y update

ENV \
    GOPATH=/go

RUN mkdir -p /go/src/github.com/ctron
ADD . /go/src/github.com/ctron/iot-simulator-operator

RUN dnf -y install golang && \
    go version && \
    cd /go/src/github.com/ctron/iot-simulator-operator/cmd/manager && go build -o /iot-simulator-operator . && \
    cd / && \
    rm -Rf go && \
    dnf -y history undo last && dnf -y clean all && \
    true

ENV IMAGE_TAG=":latest"

ENTRYPOINT /iot-simulator-operator
