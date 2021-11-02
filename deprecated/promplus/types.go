package promplus

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

// Common labels used.
const (
	LabelService = "service"
	LabelPath    = "path"
	LabelMethod  = "method"
	LabelStatus  = "status"
)

// Target describes the service that is exposing the metrics.
type Target struct {
	Namespace string
	Service   string

	durationOnce   sync.Once
	durationMetric *prometheus.HistogramVec

	visitCountOnce   sync.Once
	visitCountMetric *prometheus.CounterVec

	statusCountOnce   sync.Once
	statusCountMetric *prometheus.CounterVec
}
