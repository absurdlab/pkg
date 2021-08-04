package errorplus

import "github.com/absurdlab/pkg/timeplus"

// APIError is a standard response structure for errors. Mind you, APIError is not an implementation of error.
type APIError struct {
	Status    int                 `json:"status"`
	Timestamp *timeplus.Timestamp `json:"timestamp,omitempty"`
	RequestID string              `json:"request_id,omitempty"`
	TraceID   string              `json:"trace_id,omitempty"`
	Error     string              `json:"error"`
	Message   string              `json:"message,omitempty"`
	Detail    detail              `json:"detail,omitempty"`
	Stack     []detail            `json:"stack,omitempty"`
}
