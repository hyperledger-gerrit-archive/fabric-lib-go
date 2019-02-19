#!/bin/bash -e
#
# SPDX-License-Identifier: Apache-2.0
#

packages=(
    "github.com/client9/misspell/cmd/misspell@v0.3.4"
    "golang.org/x/tools/cmd/goimports@release-branch.go1.11"
)

for pkg in ${packages[*]}; do
    project="$(echo "$pkg" | cut -f1 -d@)"
    version="$(echo "$pkg" | cut -f2 -d@)"

    go get -d -u "${project}"
    git -C "/go/src/${project}" checkout "${version}"

    go install "${project}"
done
