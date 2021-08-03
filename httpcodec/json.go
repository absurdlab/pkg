package httpcodec

import (
	"encoding/json"
	"io"
	"net/http"
)

// EncodeJSON is an Encoder that encodes payload as JSON.
func EncodeJSON(writer io.Writer, payload interface{}) error {
	return json.NewEncoder(writer).Encode(payload)
}

// DecodeJSON produces a Decoder that decodes response body into given destination.
// The response body is closed afterwards.
func DecodeJSON(destination interface{}) Decoder {
	return func(resp *http.Response) error {
		defer func() {
			_ = resp.Body.Close()
		}()
		return json.NewDecoder(resp.Body).Decode(destination)
	}
}
