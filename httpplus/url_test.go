package httpplus_test

import (
	"github.com/absurdlab/pkg/httpplus"
	"testing"
)

func TestURLValues(t *testing.T) {
	v := httpplus.URLValues(
		"Content-Type", "application/json",
		"Accept", "application/json",
	)

	if len(v) != 2 {
		t.Error("expect url.Values to have 2 entries")
	}

	if contentType := v.Get("Content-Type"); len(contentType) == 0 {
		t.Error("expect url.Values to have Content-Type")
	}

	if accept := v.Get("Accept"); len(accept) == 0 {
		t.Error("expect url.Values to have Accept")
	}
}
