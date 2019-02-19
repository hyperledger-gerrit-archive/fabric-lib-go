# Copyright Greg Haskins All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

ARG GO_VER=1.11.5
ARG ALPINE_VER=3.8

FROM golang:${GO_VER}-alpine${ALPINE_VER} AS base
RUN apk add --no-cache \
        gcc \
        musl-dev \
        bash \
        make \
        git;

RUN mkdir -p $GOPATH/src/github.com/hyperledger/fabric-lib-go
COPY gotools.sh/ $GOPATH/src/github.com/hyperledger/fabric-lib-go
WORKDIR $GOPATH/src/github.com/hyperledger/fabric-lib-go
RUN bash -c "./gotools.sh"

RUN cp -r /go/bin /usr/local/bin

WORKDIR ${GOPATH}/src/github.com/hyperledger/fabric-lib-go
