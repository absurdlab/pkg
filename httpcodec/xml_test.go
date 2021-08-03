package httpcodec_test

import (
	"encoding/xml"
	"github.com/absurdlab/pkg/httpcodec"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestEncodeXML(t *testing.T) {
	var sb strings.Builder

	err := httpcodec.EncodeXML(&sb, xmlPayload{Message: "hello"})
	if err != nil {
		t.Error(err)
	}

	if sb.String() != "<greeting message=\"hello\"></greeting>" {
		t.Error("xml encode mismatch with expectation")
	}
}

func TestDecodeXML(t *testing.T) {
	var destination xmlPayload

	err := httpcodec.DecodeXML(&destination)(&http.Response{
		Body: io.NopCloser(strings.NewReader("<greeting message=\"hello\"></greeting>")),
	})
	if err != nil {
		t.Error(err)
	}

	if destination.Message != "hello" {
		t.Error("xml decoding error")
	}
}

type xmlPayload struct {
	XMLName xml.Name `xml:"greeting"`
	Message string   `xml:"message,attr"`
}
