package otelplus_test

import (
	"context"
	"github.com/absurdlab/pkg/otelplus"
	"testing"
)

func TestInit(t *testing.T) {
	_, closer, err := otelplus.Init(context.Background(),
		otelplus.Options.ParentOrRatioBasedSample(0.6),
		otelplus.Options.ServiceResource("test"),
		otelplus.Options.StdoutExporter(true),
	)

	if closer != nil {
		defer func() {
			_ = closer()
		}()
	}

	if err != nil {
		t.Error(err)
	}
}
