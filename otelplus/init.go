package otelplus

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Closer should be invoked to properly shut down the trace provider.
type Closer func() error

// Init returns a new trace provider initialized according to the given options and sets it as the global provider.
// It also returns a shutdown hook to be invoked to properly shut down the trace provider.
func Init(ctx context.Context, options ...TracerProviderOptionFunc) (*sdktrace.TracerProvider, Closer, error) {
	var opts []sdktrace.TracerProviderOption
	for _, each := range options {
		opt, err := each()
		if err != nil {
			return nil, nil, err
		}
		opts = append(opts, opt)
	}

	tp := sdktrace.NewTracerProvider(opts...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, func() error {
		return tp.Shutdown(ctx)
	}, nil
}
