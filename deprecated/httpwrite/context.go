package httpwrite

import (
	"github.com/absurdlab/pkg/httpcodec"
	"net/http"
)

// Options is the entrypoint for configuration WriteOptions. By default, it assumes 200 status and plain body
func Options() *WriteOptions {
	return &WriteOptions{
		status:  http.StatusOK,
		headers: map[string]string{},
		encoder: httpcodec.EncodeRaw,
	}
}

type WriteOptions struct {
	status  int
	headers map[string]string
	body    interface{}
	encoder httpcodec.Encoder
}

func (s *WriteOptions) sanitize() {
	if s.status < 100 {
		s.status = http.StatusOK
	}

	if s.encoder == nil && s.body != nil {
		s.encoder = httpcodec.EncodeRaw
	}
}

// WithStatus sets the response status.
func (s *WriteOptions) WithStatus(status int) *WriteOptions {
	s.status = status
	return s
}

// AddHeaders adds the key value pairs to response headers. If the pairs are not even numbered, the method panics.
func (s *WriteOptions) AddHeaders(kvs ...string) *WriteOptions {
	if len(kvs)%2 != 0 {
		panic("kvs must be supplied in pairs")
	}
	for i := 0; i < len(kvs); i += 2 {
		s.headers[kvs[i]] = kvs[i+1]
	}
	return s
}

// JSON sets the payload and also the "Content-Type" header to "application/json". If payload is nil, this method is noop.
func (s *WriteOptions) JSON(payload interface{}) *WriteOptions {
	if payload != nil {
		s.WithBody(payload, httpcodec.EncodeJSON).
			AddHeaders("Content-Type", "application/json")
	}
	return s
}

// XML sets the payload and also the "Content-Type" header to "application/xml". If payload is nil, this method is noop.
func (s *WriteOptions) XML(payload interface{}) *WriteOptions {
	if payload != nil {
		s.WithBody(payload, httpcodec.EncodeXML).
			AddHeaders("Content-Type", "application/xml")
	}
	return s
}

// Form sets the payload and also the "Content-Type" header to "application/x-www-form-urlencoded". If payload is nil, this method is noop.
// Check github.com/absurdlab/pkg/httpcodec for accepted types.
func (s *WriteOptions) Form(payload interface{}) *WriteOptions {
	if payload != nil {
		s.WithBody(payload, httpcodec.EncodeForm).
			AddHeaders("Content-Type", "application/x-www-form-urlencoded")
	}
	return s
}

// PlainText sets the payload and also the "Content-Type" header to "text/plain". If payload is nil, this method is noop.
// Check github.com/absurdlab/pkg/httpcodec for accepted types.
func (s *WriteOptions) PlainText(payload interface{}) *WriteOptions {
	if payload != nil {
		s.WithBody(payload, httpcodec.EncodeRaw).
			AddHeaders("Content-Type", "text/plain")
	}
	return s
}

// HTML sets the payload and also the "Content-Type" header to "text/html". If payload is nil, this method is noop.
// Check github.com/absurdlab/pkg/httpcodec for accepted types.
func (s *WriteOptions) HTML(payload interface{}) *WriteOptions {
	if payload != nil {
		s.WithBody(payload, httpcodec.EncodeRaw).
			AddHeaders("Content-Type", "text/html")
	}
	return s
}

// WithBody sets custom payload and encoder. If either is nil, this method is noop.
func (s *WriteOptions) WithBody(body interface{}, codec httpcodec.Encoder) *WriteOptions {
	if body != nil && codec != nil {
		s.body = body
		s.encoder = codec
	}
	return s
}
