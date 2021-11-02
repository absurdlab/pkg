package httpcodec_test

import (
	"github.com/absurdlab/pkg/httpcodec"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestEncodeRaw(t *testing.T) {
	var sb strings.Builder

	err := httpcodec.EncodeRaw(&sb, "hello")
	if err != nil {
		t.Error(err)
	}

	if sb.String() != "hello" {
		t.Error("raw encode mismatch")
	}
}

func TestDecodeRaw(t *testing.T) {
	var sb strings.Builder

	err := httpcodec.DecodeRaw(&sb)(&http.Response{
		Body: io.NopCloser(strings.NewReader("hello")),
	})
	if err != nil {
		t.Error(err)
	}

	if sb.String() != "hello" {
		t.Error("raw decode mismatch")
	}
}
