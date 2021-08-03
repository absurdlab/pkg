package httpcodec

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// EncodeForm encodes the payload as URL encoded form into writer. Accepted types include
// url.Values, map[string][]string, and map[string]string. Other types will error.
func EncodeForm(writer io.Writer, payload interface{}) error {
	if values, ok := payload.(url.Values); ok {
		_, err := writer.Write([]byte(values.Encode()))
		return err
	}

	switch payload.(type) {
	case map[string][]string:
		return EncodeForm(writer, url.Values(payload.(map[string][]string)))
	case map[string]string:
		values := url.Values{}
		for k, v := range payload.(map[string]string) {
			values.Add(k, v)
		}
		return EncodeForm(writer, values)
	default:
		return errors.New("invalid type for form encoder")
	}
}

// DecodeForm decodes response body as URL encoded form into the destination url.Values. The response body
// is closed afterwards.
func DecodeForm(destination *url.Values) Decoder {
	return func(resp *http.Response) error {
		defer func() {
			_ = resp.Body.Close()
		}()

		var sb strings.Builder
		if _, err := io.Copy(&sb, resp.Body); err != nil {
			return err
		}

		v, err := url.ParseQuery(sb.String())
		if err != nil {
			return err
		}

		*destination = v

		return nil
	}
}
