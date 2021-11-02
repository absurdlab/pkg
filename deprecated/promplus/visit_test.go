package promplus_test

import (
	"github.com/absurdlab/pkg/promplus"
	"github.com/absurdlab/pkg/promplus/internal"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDurationMetric(t *testing.T) {
	target := &promplus.Target{
		Namespace: "promplus",
		Service:   "test",
	}

	var h http.Handler = http.HandlerFunc(internal.Handle)
	h = target.Duration()(h)

	r := httptest.NewRequest(http.MethodGet, "/test?duration=1", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, r)

	if got, want := testutil.CollectAndCount(target.DurationMetric()), 1; got != want {
		t.Errorf("expect %d, get %d", want, got)
	}
}

func TestStatusCountMetric(t *testing.T) {
	target := &promplus.Target{
		Namespace: "promplus",
		Service:   "test",
	}

	var h http.Handler = http.HandlerFunc(internal.Handle)
	h = target.StatusCount()(h)

	r1 := httptest.NewRequest(http.MethodGet, "/test?status=200", nil)
	r2 := httptest.NewRequest(http.MethodGet, "/test?status=400", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, r1)
	h.ServeHTTP(rr, r2)
	h.ServeHTTP(rr, r1)

	expected := `
		# HELP promplus_test_http_response_status Status of HTTP response
        # TYPE promplus_test_http_response_status counter
        promplus_test_http_response_status{method="GET",path="/test",service="test",status="200"} 2
        promplus_test_http_response_status{method="GET",path="/test",service="test",status="400"} 1
`

	err := testutil.CollectAndCompare(target.StatusCountMetric(), strings.NewReader(expected))
	if err != nil {
		t.Error(err)
	}
}

func TestVisitCountMetric(t *testing.T) {
	target := &promplus.Target{
		Namespace: "promplus",
		Service:   "test",
	}

	var h http.Handler = http.HandlerFunc(internal.Handle)
	h = target.VisitCount()(h)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, r)
	h.ServeHTTP(rr, r)
	h.ServeHTTP(rr, r)

	expected := `
		# HELP promplus_test_http_requests_total Total number of requests
		# TYPE promplus_test_http_requests_total counter
		promplus_test_http_requests_total{method="GET",path="/test",service="test"} 3
	`

	err := testutil.CollectAndCompare(target.VisitCountMetric(), strings.NewReader(expected))
	if err != nil {
		t.Error(err)
	}
}
