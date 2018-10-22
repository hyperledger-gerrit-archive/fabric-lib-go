/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package healthz

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestWriteHTTPResponse(t *testing.T) {
	t.Parallel()

	now := time.Now()
	var tests = []struct {
		name         string
		hs           HealthStatus
		expectedCode int
	}{
		{
			name: "Status OK",
			hs: HealthStatus{
				Status: StatusOK,
				Time:   now,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Single Failed Check",
			hs: HealthStatus{
				Status: StatusUnavailable,
				Time:   now,
				FailedChecks: []FailedCheck{
					{
						Component: "component1",
						Reason:    "poorly written code",
					},
				},
			},
			expectedCode: http.StatusServiceUnavailable,
		},
		{
			name: "Multiple Failed Checks",
			hs: HealthStatus{
				Status: StatusUnavailable,
				Time:   now,
				FailedChecks: []FailedCheck{
					{
						Component: "component1",
						Reason:    "poorly written code",
					},
					{
						Component: "component2",
						Reason:    "more poorly written code",
					},
				},
			},
			expectedCode: http.StatusServiceUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)
			expectedJSON, _ := json.Marshal(test.hs)
			recorder := httptest.NewRecorder()
			writeHTTPResponse(recorder, test.hs)
			res := recorder.Result()
			g.Expect(res.StatusCode).To(Equal(test.expectedCode))
			defer res.Body.Close()
			actualJSON, _ := ioutil.ReadAll(res.Body)
			g.Expect(actualJSON).To(MatchJSON(expectedJSON))
		})
	}

}
