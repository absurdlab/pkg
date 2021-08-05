package otelplus

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

// DefaultHTTPClient is otelhttp.DefaultClient
var DefaultHTTPClient = otelhttp.DefaultClient

// AutoTrace returns a new HTTP middleware that automatically starts a trace, or inherits an existing trace.
func AutoTrace(operation string, options ...otelhttp.Option) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return otelhttp.NewHandler(handler, operation, options...)
	}
}
