package promplus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
	"strings"
)

// StatusCount returns an HTTP handler that counts the occurrence of status code per path, method and status.
// Multiple calls to StatusCount will only initialize the metric once.
func (t *Target) StatusCount() func(handler http.Handler) http.Handler {
	t.statusCountOnce.Do(func() {
		serviceName := strings.ToLower(t.Service)

		t.statusCountMetric = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: t.Namespace,
				Subsystem: serviceName,
				Name:      "http_response_status",
				Help:      "Status of HTTP response",
				ConstLabels: prometheus.Labels{
					LabelService: serviceName,
				},
			},
			[]string{LabelPath, LabelMethod, LabelStatus},
		)
	})

	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			delegate := &responseWriteDelegate{ResponseWriter: rw, statusCode: http.StatusOK}

			handler.ServeHTTP(delegate, r)

			t.statusCountMetric.With(prometheus.Labels{
				LabelPath:   r.URL.Path,
				LabelMethod: r.Method,
				LabelStatus: strconv.Itoa(delegate.statusCode),
			}).Inc()
		})
	}
}

// StatusCountMetric exposes the status count metric. It may or may not have been initialized.
func (t *Target) StatusCountMetric() prometheus.Collector {
	return t.statusCountMetric
}

type responseWriteDelegate struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriteDelegate) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
