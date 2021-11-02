package promplus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strings"
)

// Duration returns an HTTP middleware that records the request duration per HTTP method and path. Multiple
// calls of Duration on the same Target will only register the metric once.
func (t *Target) Duration() func(handler http.Handler) http.Handler {
	t.durationOnce.Do(func() {
		serviceName := strings.ToLower(t.Service)

		t.durationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: t.Namespace,
			Subsystem: serviceName,
			Name:      "http_response_time_seconds",
			Help:      "Duration of http handler response time in seconds",
			ConstLabels: prometheus.Labels{
				LabelService: serviceName,
			},
		}, []string{LabelPath, LabelMethod})
	})

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			timer := prometheus.NewTimer(t.durationMetric.With(prometheus.Labels{
				LabelPath:   r.URL.Path,
				LabelMethod: r.Method,
			}))
			defer timer.ObserveDuration()

			handler.ServeHTTP(rw, r)
		})
	}
}

// DurationMetric exposes the internal duration metric. It may or may not have been initialized.
func (t *Target) DurationMetric() prometheus.Collector {
	return t.durationMetric
}
