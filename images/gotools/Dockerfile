# Copyright Greg Haskins All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

ARG GO_VER=1.11.5
ARG ALPINE_VER=3.8

FROM golang:${GO_VER}-alpine${ALPINE_VER} AS base
ARG PROJECT
RUN apk add --no-cache \
        gcc \
        musl-dev \
        bash \
        git;

ADD . $GOPATH/src/github.com/hyperledger/${PROJECT}
WORKDIR $GOPATH/src/github.com/hyperledger/${PROJECT}
RUN bash -c "./images/gotools/gotools.sh"

FROM golang:${GO_VER}-alpine${ALPINE_VER} AS gotools
ARG PROJECT
RUN apk add --no-cache \
        gcc \
        musl-dev \
        bash \
        make \
        git;

COPY --from=base /go/bin/misspell /usr/local/bin
COPY --from=base /go/bin/goimports /usr/local/bin

WORKDIR ${GOPATH}/src/github.com/hyperledger/${PROJECT}
