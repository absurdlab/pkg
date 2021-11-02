package promplus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strings"
)

// VisitCount returns an HTTP middleware that keeps track of the number of times the endpoint is
// visited per path and method. Multiple calls to VisitCount will only initialize the metric once.
func (t *Target) VisitCount() func(handler http.Handler) http.Handler {
	t.visitCountOnce.Do(func() {
		serviceName := strings.ToLower(t.Service)

		t.visitCountMetric = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: t.Namespace,
			Subsystem: serviceName,
			Name:      "http_requests_total",
			Help:      "Total number of requests",
			ConstLabels: prometheus.Labels{
				LabelService: serviceName,
			},
		}, []string{LabelPath, LabelMethod})
	})

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			t.visitCountMetric.With(prometheus.Labels{
				LabelPath:   r.URL.Path,
				LabelMethod: r.Method,
			}).Inc()

			handler.ServeHTTP(rw, r)
		})
	}
}

// VisitCountMetric exposes the visit count metric. It may or may not have been initialized.
func (t *Target) VisitCountMetric() prometheus.Collector {
	return t.visitCountMetric
}
