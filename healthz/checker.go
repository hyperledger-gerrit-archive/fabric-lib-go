/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

// Package healthz provides an HTTP handler which returns the health status of
// one or more components of an application or service.
package healthz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type AlreadyRegisteredError string

func (are AlreadyRegisteredError) Error() string {
	return fmt.Sprintf("'%s' is already registered", are)
}

const (
	// StatusOK is returned if all health checks pass.
	StatusOK = "OK"
	// StatusUnavailable is returned if any health check fails.
	StatusUnavailable = "Service Unavailable"
)

//go:generate counterfeiter -o mock/health_checker.go -fake-name HealthChecker . HealthChecker

// HealthChecker defines the interface components must implement in order to
// register with the Handler in order to be included in the health status.
type HealthChecker interface {
	HealthCheck() (message string, ok bool)
}

// FailedCheck represents a failed status check for a component.
type FailedCheck struct {
	Component string `json:"component"`
	Reason    string `json:"reason"`
}

// HealthStatus represents the current health status of all registered components.
type HealthStatus struct {
	Status       string        `json:"status"`
	Time         time.Time     `json:"time"`
	FailedChecks []FailedCheck `json:"failed_checks,omitempty"`
}

// NewHealthHandler returns a new HealthHandler instance.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		HealthCheckers: map[string]HealthChecker{},
	}
}

// HealthHandler is responsible for executing registered health checks.  It
// provides an HTTP handler which returns the health status for all registered
// components.
type HealthHandler struct {
	mutex          sync.Mutex
	HealthCheckers map[string]HealthChecker
}

// RegisterChecker registers a HealthChecker for a named component and adds it to
// the list of status checks to run.  It returns an error if the component has
// already been registered.
func (h *HealthHandler) RegisterChecker(component string, checker HealthChecker) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.HealthCheckers[component]; ok {
		return AlreadyRegisteredError(component)
	}
	h.HealthCheckers[component] = checker
	return nil
}

// DeregisterChecker degisters a named HealthChecker.
func (h *HealthHandler) DeregisterChecker(component string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.HealthCheckers[component]; !ok {
		return
	}
	delete(h.HealthCheckers, component)
}

// HandleHTTP returns an HTTP handler (see http.Handler) which can be used as
// an HTTP endpoint for health checks.  If all registered checks pass, it returns
// an HTTP status `200 OK` with a JSON payload of `{"status": "OK"}`.  If all
// checks do not pass, it returns an HTTP status `503 Service Unavailable` with
// a JSON payload of `{"status": "OK","failed_checks":[...]}.
func (h *HealthHandler) HandleHTTP() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		hs := HealthStatus{
			Status: StatusOK,
			Time:   time.Now(),
		}

		failedChecks := h.RunChecks()
		if len(failedChecks) > 0 {
			hs.Status = StatusUnavailable
			hs.FailedChecks = failedChecks
		}
		writeHTTPResponse(rw, hs)
	})
}

// RunChecks runs all HealthCheckers and returns any failures.
func (h *HealthHandler) RunChecks() []FailedCheck {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var failedChecks []FailedCheck

	for component, checker := range h.HealthCheckers {
		if reason, ok := checker.HealthCheck(); !ok {
			failedCheck := FailedCheck{
				Component: component,
				Reason:    reason,
			}
			failedChecks = append(failedChecks, failedCheck)
		}
	}
	return failedChecks
}

// write the HTTP response
func writeHTTPResponse(rw http.ResponseWriter, hs HealthStatus) {
	var resp []byte
	rc := http.StatusOK
	if len(hs.FailedChecks) > 0 {
		rc = http.StatusServiceUnavailable
	}
	// NOTE: no error check below as it seems unnecessary
	resp, _ = json.Marshal(hs)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(rc)
	rw.Write(resp)
}
