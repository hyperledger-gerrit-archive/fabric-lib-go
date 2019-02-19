#!/bin/bash -e

# Copyright Greg Haskins All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

MISSPELL=v0.3.4
GOIMPORTS=release-branch.go1.11

go get -d -u github.com/client9/misspell/cmd/misspell
echo "git -C ${GOPATH}/src/github.com/client9/misspell/cmd/misspell checkout ${MISSPELL}"
go install github.com/client9/misspell/cmd/misspell

go get -d -u golang.org/x/tools/cmd/goimports
echo "git -C /go/src/golang.org/x/tools/cmd/goimports checkout -q ${GOIMPORTS}"
go install golang.org/x/tools/cmd/goimports
