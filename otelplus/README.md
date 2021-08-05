# otelplus

Opinionated convenience for setting up and using the OpenTelemetry tracing library.

```shell
go get -u github.com/absurdlab/pkg/otelplus
```

## Usage

```go
// Setup a trace provider
provider, closer, _ := otelplus.Init(context.Background(),
    otelplus.Options.ParentOrRatioBasedSample(0.6),
    otelplus.Options.ServiceResource("my_app"),
    otelplus.Options.StdoutExporter(true),
)
defer func() {
    _ = closer()
}

// Setup an HTTP middleware to automatically trace
h := myHandler()
h = otelplus.AutoTrace("get_user")(h)

// Get trace id and span id from context
traceID := otelplus.TraceID(ctx)
spanID := otelplus.SpanID(ctx)

// Add details to trace
otelplus.AddEvent(ctx, "bang", map[string]interface{}{
    "foo": "bar",
    "count": 3,
})

// Send and receive baggage, use context to propagate and receive
ctx = otelplus.ContextWithBaggage(ctx, "foo", "bar") // propagate
_ = BagMemberValue(ctx, "foo") // receive
```