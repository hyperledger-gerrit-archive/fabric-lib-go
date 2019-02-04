# Copyright IBM Corp All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# -------------------------------------------------------------
# This makefile defines the following targets

#   - all (default) - runs checks and unit-tests
#   - unit-test - runs the go-test based unit tests
#   - checks - runs all check conditions (spellcheck, license, trailing-spaces and vet)

.PHONY: all

all: checks unit-tests

checks: spellcheck license trailing-spaces vet

PHONY: unit-tests
unit-tests:
	go test -race ./... -short -p 1 -timeout=20m

.PHONY: spellcheck
spellcheck:
	@scripts/check_spelling

.PHONY: license
license:
	@scripts/check_license

.PHONY: trailing-spaces
trailing-spaces:
	@scripts/check_trailingspaces

.PHONY: vet
vet:
	@echo "Running vet checks.."
	@scripts/check_vet
