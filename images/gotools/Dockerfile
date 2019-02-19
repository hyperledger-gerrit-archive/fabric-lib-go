ARG GO_VER=1.11.5
ARG ALPINE_VER=3.8
ARG MISSPELL=c0b55c8
ARG GOIMPORTS=8dcb7bc

FROM golang:${GO_VER}-alpine${ALPINE_VER} as base

RUN apk add --no-cache \
	gcc \
	musl-dev \
        bash \
        make \
	git

RUN go get -d -u github.com/client9/misspell/cmd/misspell
RUN git -C /go/src/github.com/client9/misspell/cmd/misspell checkout -q ${MISSPELL}
RUN go install github.com/client9/misspell/cmd/misspell
RUN git -C /go/src/github.com/client9/misspell/cmd/misspell checkout -q master

RUN go get -d -u golang.org/x/tools/cmd/goimports
RUN git -C /go/src/golang.org/x/tools/cmd/goimports checkout -q ${GOIMPORTS}
RUN go install golang.org/x/tools/cmd/goimports
RUN git -C /go/src/golang.org/x/tools/cmd/goimports checkout -q master


FROM golang:${GO_VER}-alpine${ALPINE_VER} as gotools

RUN apk add --no-cache \
        gcc \
        musl-dev \
        bash \
        make \
        git

COPY --from=base /go/bin/misspell /usr/local/bin
COPY --from=base /go/bin/goimports /usr/local/bin

WORKDIR /go/src/github.com/hyperledger/fabric-lib-go
