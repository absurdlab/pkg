package httpcodec

import (
	"io"
	"net/http"
)

// Encoder encodes payload into writer.
type Encoder func(writer io.Writer, payload interface{}) error

// Decoder decodes response body into some destination of the caller's choice.
type Decoder func(resp *http.Response) error
