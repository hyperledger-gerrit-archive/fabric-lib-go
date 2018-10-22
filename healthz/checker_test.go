/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package healthz_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyperledger/fabric-lib-go/healthz"
	"github.com/hyperledger/fabric-lib-go/healthz/mock"
	. "github.com/onsi/gomega"
)

func TestRegisterChecker(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	handler := healthz.NewHealthHandler()
	checker1 := &mock.HealthChecker{}
	checker2 := &mock.HealthChecker{}
	component1 := "component1"
	component2 := "component2"

	err := handler.RegisterChecker(component1, checker1)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(handler.HealthCheckers).To(HaveKey(component1))

	err = handler.RegisterChecker(component1, checker1)
	g.Expect(err).To(MatchError(healthz.AlreadyRegisteredError(component1)))

	err = handler.RegisterChecker(component2, checker2)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(handler.HealthCheckers).To(HaveLen(2))
	g.Expect(handler.HealthCheckers).To(HaveKey(component1))
	g.Expect(handler.HealthCheckers).To(HaveKey(component2))

	handler.DeregisterChecker(component1)
	g.Expect(handler.HealthCheckers).To(HaveLen(1))
	g.Expect(handler.HealthCheckers).ToNot(HaveKey(component1))
	g.Expect(handler.HealthCheckers).To(HaveKey(component2))

	// deregister non-existent checker should be a no-op
	handler.DeregisterChecker(component1)
	g.Expect(handler.HealthCheckers).To(HaveLen(1))
}

func TestRunChecks(t *testing.T) {
	t.Parallel()
	g := NewGomegaWithT(t)

	handler := healthz.NewHealthHandler()

	checker := &mock.HealthChecker{}
	checker.HealthCheckReturns("", true)
	handler.RegisterChecker("good_component", checker)
	fc := handler.RunChecks()
	g.Expect(fc).To(HaveLen(0))

	failedChecker := &mock.HealthChecker{}
	reason := "poorly written code"
	failedChecker.HealthCheckReturnsOnCall(0, reason, false)
	failedChecker.HealthCheckReturnsOnCall(1, reason, false)
	handler.RegisterChecker("bad_component1", failedChecker)
	handler.RegisterChecker("bad_component2", failedChecker)
	fc = handler.RunChecks()
	g.Expect(failedChecker.HealthCheckCallCount()).To(Equal(2))
	g.Expect(fc).To(HaveLen(2))
	g.Expect(fc).To(ContainElement(
		healthz.FailedCheck{
			Component: "bad_component1",
			Reason:    reason,
		},
	))
	g.Expect(fc).To(ContainElement(
		healthz.FailedCheck{
			Component: "bad_component2",
			Reason:    reason,
		},
	))
}

func TestHandleHTTP(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name           string
		hh             *healthz.HealthHandler
		failedChecks   []healthz.FailedCheck
		expectedCode   int
		expectedStatus string
	}{
		{
			name: "Status OK",
			hh: &healthz.HealthHandler{
				HealthCheckers: map[string]healthz.HealthChecker{
					"component1": &mock.HealthChecker{
						HealthCheckStub: func() (string, bool) {
							return "", true
						},
					},
				},
			},
			expectedCode:   http.StatusOK,
			expectedStatus: healthz.StatusOK,
		},
		{
			name: "Service Unavailable",
			hh: &healthz.HealthHandler{
				HealthCheckers: map[string]healthz.HealthChecker{
					"component1": &mock.HealthChecker{
						HealthCheckStub: func() (string, bool) {
							return "poorly written code", false
						},
					},
				},
			},
			failedChecks: []healthz.FailedCheck{
				{
					Component: "component1",
					Reason:    "poorly written code",
				},
			},
			expectedCode:   http.StatusServiceUnavailable,
			expectedStatus: healthz.StatusUnavailable,
		},
		{
			name: "Service Unavailable - Multiple",
			hh: &healthz.HealthHandler{
				HealthCheckers: map[string]healthz.HealthChecker{
					"component1": &mock.HealthChecker{
						HealthCheckStub: func() (string, bool) {
							return "poorly written code", false
						},
					},
					"component2": &mock.HealthChecker{
						HealthCheckStub: func() (string, bool) {
							return "more poorly written code", false
						},
					},
				},
			},
			failedChecks: []healthz.FailedCheck{
				{
					Component: "component1",
					Reason:    "poorly written code",
				},
				{
					Component: "component2",
					Reason:    "more poorly written code",
				},
			},
			expectedCode:   http.StatusServiceUnavailable,
			expectedStatus: healthz.StatusUnavailable,
		},
		{
			name: "Mixed",
			hh: &healthz.HealthHandler{
				HealthCheckers: map[string]healthz.HealthChecker{
					"component1": &mock.HealthChecker{
						HealthCheckStub: func() (string, bool) {
							return "poorly written code", false
						},
					},
					"component2": &mock.HealthChecker{
						HealthCheckStub: func() (string, bool) {
							return "", true
						},
					},
				},
			},
			failedChecks: []healthz.FailedCheck{
				{
					Component: "component1",
					Reason:    "poorly written code",
				},
			},
			expectedCode:   http.StatusServiceUnavailable,
			expectedStatus: healthz.StatusUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)
			ts := httptest.NewServer(test.hh.HandleHTTP())
			defer ts.Close()

			res, err := http.Get(ts.URL)
			if err != nil {
				t.Fatalf("unexpected error getting response [%s]", err)
			}
			g.Expect(res.StatusCode).To(Equal(test.expectedCode))
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatalf("unexpected error reading response body [%s]", err)
			}
			var hs healthz.HealthStatus
			err = json.Unmarshal(body, &hs)
			if err != nil {
				t.Fatalf("unexpected error marshaling response body [%s]", err)
			}
			g.Expect(hs.Status).To(Equal(test.expectedStatus))
			for _, failedCheck := range test.failedChecks {
				g.Expect(hs.FailedChecks).To(ContainElement(failedCheck))
			}
		})
	}
}
