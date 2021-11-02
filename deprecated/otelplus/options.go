package otelplus

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

// Options is the entry point of creating options to configure the trace provider.
var Options = options{}

// TracerProviderOptionFunc can return an tracer provider option, and potentially an error.
type TracerProviderOptionFunc func() (sdktrace.TracerProviderOption, error)

type options struct{}

// AlwaysSample will configure the trace provider to sample all traces.
func (_ options) AlwaysSample() TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithSampler(sdktrace.AlwaysSample()), nil
	}
}

// NeverSample will configure the trace provider to sample no traces.
func (_ options) NeverSample() TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithSampler(sdktrace.NeverSample()), nil
	}
}

// RatioBasedSample will configure the trace provider to sample a portion of the traces.
func (_ options) RatioBasedSample(fraction float64) TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithSampler(sdktrace.TraceIDRatioBased(fraction)), nil
	}
}

// ParentOrAlwaysSample will configure the trace provider to sample based on the current span config.
// If no config was specified on current span, defaults to always sampling.
func (_ options) ParentOrAlwaysSample() TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())), nil
	}
}

// ParentOrNeverSample will configure the trace provider to sample based on the current span config.
// If no config was specified on current span, defaults to never sampling.
func (_ options) ParentOrNeverSample() TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.NeverSample())), nil
	}
}

// ParentOrRatioBasedSample will configure the trace provider to sample based on the current span config.
// If no config was specified on current span, defaults to ratio based sampling.
func (_ options) ParentOrRatioBasedSample(fraction float64) TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(fraction))), nil
	}
}

// ServiceResource will configure the trace provider to include the resource name and schema url in the trace.
func (_ options) ServiceResource(name string) TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(name),
			),
		), nil
	}
}

// InsecureGRPCExporter will configure the trace provider to export traces to a collector using OpenTelemetry gRPC
// API over an insecure connection.
func (_ options) InsecureGRPCExporter(ctx context.Context, url string) TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		exp, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(url),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		)
		if err != nil {
			return nil, err
		}

		return sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exp)), nil
	}
}

// StdoutExporter will configure the trace provider to export to stdout.
func (_ options) StdoutExporter(prettyPrint bool) TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		var opt []stdouttrace.Option
		if prettyPrint {
			opt = append(opt, stdouttrace.WithPrettyPrint())
		}

		exp, err := stdouttrace.New(opt...)
		if err != nil {
			return nil, err
		}

		return sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exp)), nil
	}
}

// NoopExporter will configure the trace provider to do nothing during export.
func (_ options) NoopExporter() TracerProviderOptionFunc {
	return func() (sdktrace.TracerProviderOption, error) {
		return sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(noopExporter{})), nil
	}
}

type noopExporter struct{}

func (_ noopExporter) ExportSpans(_ context.Context, _ []sdktrace.ReadOnlySpan) error {
	return nil
}

func (_ noopExporter) Shutdown(_ context.Context) error {
	return nil
}
