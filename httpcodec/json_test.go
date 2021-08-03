package httpcodec_test

import (
	"github.com/absurdlab/pkg/httpcodec"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestEncodeJSON(t *testing.T) {
	var sb strings.Builder

	err := httpcodec.EncodeJSON(&sb, jsonPayload{Message: "hello"})
	if err != nil {
		t.Error(err)
	}

	if sb.String() != "{\"message\":\"hello\"}\n" {
		t.Error("json payload mismatch with expectation")
	}
}

func TestDecodeJSON(t *testing.T) {
	var destination jsonPayload

	err := httpcodec.DecodeJSON(&destination)(&http.Response{
		Body: io.NopCloser(strings.NewReader("{\"message\":\"hello\"}\n")),
	})
	if err != nil {
		t.Error(err)
	}

	if destination.Message != "hello" {
		t.Error("json decoding error")
	}
}

type jsonPayload struct {
	Message string `json:"message"`
}
