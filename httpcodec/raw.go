package httpcodec

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

// EncodeRaw is an Encoder that encodes payload as is. Accept payload types are
// string, []byte, *bytes.Buffer, and *strings.Builder. Any other payload types
// will error.
func EncodeRaw(rw io.Writer, payload interface{}) error {
	if rawBytes, ok := payload.([]byte); ok {
		_, err := rw.Write(rawBytes)
		return err
	}

	switch payload.(type) {
	case string:
		return EncodeRaw(rw, []byte(payload.(string)))
	case *strings.Builder:
		return EncodeRaw(rw, []byte(payload.(*strings.Builder).String()))
	case *bytes.Buffer:
		return EncodeRaw(rw, payload.(*bytes.Buffer).Bytes())
	default:
		return errors.New("invalid type for raw encoder")
	}
}

// DecodeRaw decodes the response body into the given writer. Typically, the writer
// is backed by bytes.Buffer, or strings.Builder. The response body is closed afterwards.
func DecodeRaw(writer io.Writer) Decoder {
	return func(resp *http.Response) error {
		defer func() {
			_ = resp.Body.Close()
		}()
		_, err := io.Copy(writer, resp.Body)
		return err
	}
}
