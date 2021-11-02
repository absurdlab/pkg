package httpcodec_test

import (
	"github.com/absurdlab/pkg/httpcodec"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestEncodeForm(t *testing.T) {
	var sb strings.Builder

	err := httpcodec.EncodeForm(&sb, map[string]string{"foo": "bar", "one": "two"})
	if err != nil {
		t.Error(err)
	}

	if sb.String() != "foo=bar&one=two" {
		t.Error("form encoding mismatch")
	}
}

func TestDecodeForm(t *testing.T) {
	var form url.Values

	err := httpcodec.DecodeForm(&form)(&http.Response{
		Body: io.NopCloser(strings.NewReader("foo=bar&one=two")),
	})
	if err != nil {
		t.Error(err)
	}

	if len(form) != 2 || form.Get("foo") != "bar" || form.Get("one") != "two" {
		t.Error("form decode mismatch")
	}
}
