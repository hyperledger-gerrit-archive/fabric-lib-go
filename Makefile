# Copyright IBM Corp All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets

#   - unit-test - runs the go-test based unit tests

.PHONY: all

all: unit-tests

PHONY: unit-tests
unit-tests:
	go test ./...
