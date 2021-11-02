package otelplus

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

// TraceID extracts the trace id from the current span.
func TraceID(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).TraceID().String()
}

// SpanID extracts the span id from the current span.
func SpanID(ctx context.Context) string {
	return trace.SpanContextFromContext(ctx).SpanID().String()
}

// BagMemberValue retrieves the value of the baggage member carried by the current span.
func BagMemberValue(ctx context.Context, key string) string {
	return baggage.FromContext(ctx).Member(key).Value()
}

// ContextWithBaggage returns a new context.Context with key value pairs set as baggage.
func ContextWithBaggage(ctx context.Context, kvs ...string) context.Context {
	if len(kvs)%2 != 0 {
		panic("kvs must be provided in pairs")
	}

	var pairs []string
	for i := 0; i < len(kvs); i += 2 {
		pairs = append(pairs, fmt.Sprintf("%s=%s", kvs[i], kvs[i+1]))
	}

	bag, _ := baggage.Parse(strings.Join(pairs, ","))

	return baggage.ContextWithBaggage(ctx, bag)
}

// AddEvent adds the given event, with key value parameters to the current span. The permitted types
// for the parameters are string, int, int64, float64, bool, []string and []interface{}, otherwise the
// method panics.
func AddEvent(ctx context.Context, event string, kvs map[string]interface{}) {
	var attributes []attribute.KeyValue
	for k, v := range kvs {
		var attr attribute.KeyValue
		switch v.(type) {
		case string:
			attr = attribute.Key(k).String(v.(string))
		case int:
			attr = attribute.Key(k).Int(v.(int))
		case int64:
			attr = attribute.Key(k).Int64(v.(int64))
		case float64:
			attr = attribute.Key(k).Float64(v.(float64))
		case bool:
			attr = attribute.Key(k).Bool(v.(bool))
		case []string:
			attr = attribute.Key(k).Array(v.([]string))
		case []interface{}:
			attr = attribute.Key(k).Array(v.([]interface{}))
		default:
			panic("invalid value type")
		}
		attributes = append(attributes, attr)
	}

	trace.SpanFromContext(ctx).AddEvent(event, trace.WithAttributes(attributes...))
}
