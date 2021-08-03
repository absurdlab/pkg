package httpcall

import (
	"bytes"
	"github.com/absurdlab/pkg/httpcodec"
	"io"
	"net/http"
	"net/url"
)

// Options returns a new default call context, with the default http.Client; GET as http method;
// raw encoding for payload; 2XX status code to be considered successful response.
func Options() *CallOptions {
	return &CallOptions{
		client:    http.DefaultClient,
		method:    http.MethodGet,
		params:    url.Values{},
		headers:   map[string]string{},
		payload:   nil,
		encoder:   httpcodec.EncodeRaw,
		isSuccess: defaultIsSuccess,
	}
}

type CallOptions struct {
	client         *http.Client
	method         string
	url            string
	params         url.Values
	headers        map[string]string
	payload        interface{}
	encoder        httpcodec.Encoder
	successDecoder httpcodec.Decoder
	errorDecoder   httpcodec.Decoder
	isSuccess      func(resp *http.Response) bool
}

func (s *CallOptions) sanitize() {
	if s.client == nil {
		s.client = http.DefaultClient
	}

	if len(s.method) == 0 {
		s.method = http.MethodGet
	}

	if s.encoder == nil && s.payload != nil {
		s.encoder = httpcodec.EncodeRaw
	}

	if s.isSuccess == nil {
		s.isSuccess = defaultIsSuccess
	}
}

// WithClient configures a new http.Client.
func (s *CallOptions) WithClient(client *http.Client) *CallOptions {
	if client == nil {
		return s
	}
	s.client = client
	return s
}

// GET is shortcut for WithMethod and WithURL
func (s *CallOptions) GET(url string) *CallOptions {
	return s.WithMethod(http.MethodGet).WithURL(url)
}

// POST is shortcut for WithMethod and WithURL
func (s *CallOptions) POST(url string) *CallOptions {
	return s.WithMethod(http.MethodPost).WithURL(url)
}

// PUT is shortcut for WithMethod and WithURL
func (s *CallOptions) PUT(url string) *CallOptions {
	return s.WithMethod(http.MethodPut).WithURL(url)
}

// PATCH is shortcut for WithMethod and WithURL
func (s *CallOptions) PATCH(url string) *CallOptions {
	return s.WithMethod(http.MethodPatch).WithURL(url)
}

// DELETE is shortcut for WithMethod and WithURL
func (s *CallOptions) DELETE(url string) *CallOptions {
	return s.WithMethod(http.MethodDelete).WithURL(url)
}

// WithMethod configures a new HTTP method. If invalid HTTP method is supplied, this method is noop.
func (s *CallOptions) WithMethod(method string) *CallOptions {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodHead, http.MethodOptions, http.MethodTrace, http.MethodConnect:
		s.method = method
	}
	return s
}

// WithURL sets the target url. If supplied url is empty, the method is noop.
func (s *CallOptions) WithURL(url string) *CallOptions {
	if len(url) > 0 {
		s.url = url
	}
	return s
}

// AddParams adds the key value pairs as query parameters. If the supplied key values are not in pairs, this method panics.
func (s *CallOptions) AddParams(kvs ...string) *CallOptions {
	if len(kvs)%2 != 0 {
		panic("kvs must be supplied in pairs")
	}
	for i := 0; i < len(kvs); i += 2 {
		s.params.Add(kvs[i], kvs[i+1])
	}
	return s
}

// AddHeaders adds the key value pairs as headers. If the supplied key values are not in pairs, this method panics.
func (s *CallOptions) AddHeaders(kvs ...string) *CallOptions {
	if len(kvs)%2 != 0 {
		panic("kvs must be supplied in pairs")
	}
	for i := 0; i < len(kvs); i += 2 {
		s.headers[kvs[i]] = kvs[i+1]
	}
	return s
}

// JSON sets the payload and also the "Content-Type" header to "application/json". If payload is nil, this method is noop.
func (s *CallOptions) JSON(payload interface{}) *CallOptions {
	if payload != nil {
		s.WithPayload(payload, httpcodec.EncodeJSON).
			AddHeaders("Content-Type", "application/json")
	}
	return s
}

// XML sets the payload and also the "Content-Type" header to "application/xml". If payload is nil, this method is noop.
func (s *CallOptions) XML(payload interface{}) *CallOptions {
	if payload != nil {
		s.WithPayload(payload, httpcodec.EncodeXML).
			AddHeaders("Content-Type", "application/xml")
	}
	return s
}

// Form sets the payload and also the "Content-Type" header to "application/x-www-form-urlencoded". If payload is nil, this method is noop.
// Check github.com/absurdlab/pkg/httpcodec for accepted types.
func (s *CallOptions) Form(payload interface{}) *CallOptions {
	if payload != nil {
		s.WithPayload(payload, httpcodec.EncodeForm).
			AddHeaders("Content-Type", "application/x-www-form-urlencoded")
	}
	return s
}

// Plain sets the payload and also the "Content-Type" header to "text/plain". If payload is nil, this method is noop.
// Check github.com/absurdlab/pkg/httpcodec for accepted types.
func (s *CallOptions) Plain(payload interface{}) *CallOptions {
	if payload != nil {
		s.WithPayload(payload, httpcodec.EncodeRaw).
			AddHeaders("Content-Type", "text/plain")
	}
	return s
}

// WithPayload sets custom payload and encoder. If either is nil, this method is noop.
func (s *CallOptions) WithPayload(payload interface{}, codec httpcodec.Encoder) *CallOptions {
	if payload != nil && codec != nil {
		s.payload = payload
		s.encoder = codec
	}
	return s
}

// ToJSONSuccess sets success response codec to json decoder and sets "Accept" header to "application/json". If
// destination is nil, this method is noop.
func (s *CallOptions) ToJSONSuccess(destination interface{}) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "application/json").
			ToSuccess(httpcodec.DecodeJSON(destination))
	}
	return s
}

// ToJSONError sets error response codec to json decoder and sets "Accept" header to "application/json". If
// destination is nil, this method is noop.
func (s *CallOptions) ToJSONError(destination interface{}) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "application/json").
			ToError(httpcodec.DecodeJSON(destination))
	}
	return s
}

// ToXMLSuccess sets success response codec to xml decoder and sets "Accept" header to "application/xml". If
// destination is nil, this method is noop.
func (s *CallOptions) ToXMLSuccess(destination interface{}) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "application/xml").
			ToSuccess(httpcodec.DecodeXML(destination))
	}
	return s
}

// ToXMLError sets error response codec to xml decoder and sets "Accept" header to "application/xml". If
// destination is nil, this method is noop.
func (s *CallOptions) ToXMLError(destination interface{}) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "application/xml").
			ToError(httpcodec.DecodeXML(destination))
	}
	return s
}

// ToFormSuccess sets success response codec to form decoder and sets "Accept" header to "application/x-www-form-urlencoded".
// If destination is nil, this method is noop.
func (s *CallOptions) ToFormSuccess(destination *url.Values) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "application/x-www-form-urlencoded").
			ToSuccess(httpcodec.DecodeForm(destination))
	}
	return s
}

// ToFormError sets error response codec to form decoder and sets "Accept" header to "application/x-www-form-urlencoded".
// If destination is nil, this method is noop.
func (s *CallOptions) ToFormError(destination *url.Values) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "application/x-www-form-urlencoded").
			ToError(httpcodec.DecodeForm(destination))
	}
	return s
}

// ToPlainSuccess sets success response codec to form decoder and sets "Accept" header to "text/plain".
// If destination is nil, this method is noop.
func (s *CallOptions) ToPlainSuccess(destination io.Writer) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "text/plain").
			ToSuccess(httpcodec.DecodeRaw(destination))
	}
	return s
}

// ToPlainError sets error response codec to form decoder and sets "Accept" header to "text/plain".
// If destination is nil, this method is noop.
func (s *CallOptions) ToPlainError(destination io.Writer) *CallOptions {
	if destination != nil {
		s.AddHeaders("Accept", "text/plain").
			ToError(httpcodec.DecodeRaw(destination))
	}
	return s
}

// ToSuccess sets success response codec. If codec is nil, this method is noop.
func (s *CallOptions) ToSuccess(codec httpcodec.Decoder) *CallOptions {
	if codec != nil {
		s.successDecoder = codec
	}
	return s
}

// ToError sets error response codec. If codec is nil, this method is noop.
func (s *CallOptions) ToError(codec httpcodec.Decoder) *CallOptions {
	if codec != nil {
		s.errorDecoder = codec
	}
	return s
}

// IsSuccessWhenStatus sets the success criteria to that the response is only successful when it returns one of the
// supplied status code.
func (s *CallOptions) IsSuccessWhenStatus(statuses ...int) *CallOptions {
	s.isSuccess = func(resp *http.Response) bool {
		for _, each := range statuses {
			if resp.StatusCode == each {
				return true
			}
		}
		return false
	}
	return s
}

// IsSuccessWhenStatusInRange sets the success criteria to that the response is only successful when it returns a status
// code in the supplied range.
func (s *CallOptions) IsSuccessWhenStatusInRange(lowerInclusive int, upperInclusive int) *CallOptions {
	s.isSuccess = func(resp *http.Response) bool {
		if resp.StatusCode >= lowerInclusive && resp.StatusCode <= upperInclusive {
			return true
		}
		return false
	}
	return s
}

// WithSuccessCriteria sets the success criteria. If criteria is nil, this method is noop.
func (s *CallOptions) WithSuccessCriteria(criteria func(resp *http.Response) bool) *CallOptions {
	if criteria != nil {
		s.isSuccess = criteria
	}
	return s
}

func (s *CallOptions) urlString() (string, error) {
	if len(s.params) == 0 {
		return s.url, nil
	}

	u, err := url.Parse(s.url)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, vs := range s.params {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (s *CallOptions) body() (io.Reader, error) {
	var body io.Reader
	if s.payload != nil {
		body = new(bytes.Buffer)
		if err := s.encoder(body.(*bytes.Buffer), s.payload); err != nil {
			return nil, err
		}
	}
	return body, nil
}

var defaultIsSuccess = func(resp *http.Response) bool {
	return resp.StatusCode > 199 && resp.StatusCode < 300
}
