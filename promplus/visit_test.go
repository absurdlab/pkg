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
