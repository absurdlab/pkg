package httpplus

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestSpec specifies the request specification and parsing strategies.
type RequestSpec struct {
	// Client is the http.Client to use instead of http.DefaultClient
	Client *http.Client
	// Method is the http method, defaults to GET
	Method string
	// URL is the target URL of the request.
	URL string
	// Headers is the headers to add to request.
	Headers map[string]string
	// Payload is the object to be sent as request body, it will be encoded as Encoder.
	Payload interface{}
	// Encoder is the request body encoder, defaults to JSONEncoder.
	Encoder Encoder
	// SuccessDecoder decodes response body when it is successful. Successful-ness is
	// judged by IsSuccess. If not specified, success payload will not be parsed.
	SuccessDecoder Decoder
	// ErrorDecoder decodes response body when it has failed. Failed-ness is judged
	// by IsSuccess. If not specified, error payload will not be parsed.
	ErrorDecoder Decoder
	// IsSuccess judges whether a response is successful or not. By default, all
	// response with 2XX status codes are deemed successful.
	IsSuccess IsSuccess
}

// MakeRequest executes the specification.
func MakeRequest(ctx context.Context, spec *RequestSpec) (*http.Response, error) {
	return spec.execute(ctx)
}

func (s *RequestSpec) sanitize() {
	if s.Client == nil {
		s.Client = http.DefaultClient
	}

	if len(s.Method) == 0 {
		s.Method = http.MethodGet
	}

	if s.Payload != nil && s.Encoder == nil {
		s.Encoder = JSONEncoder
	}

	if s.IsSuccess == nil {
		s.IsSuccess = DefaultIsSuccess
	}
}

func (s *RequestSpec) execute(ctx context.Context) (*http.Response, error) {
	s.sanitize()

	var body io.Reader
	if s.Payload != nil {
		body = new(bytes.Buffer)
		if err := s.Encoder(body.(*bytes.Buffer), s.Payload); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, s.Method, s.URL, body)
	if err != nil {
		return nil, err
	}

	if len(s.Headers) > 0 {
		for k, v := range s.Headers {
			request.Header.Set(k, v)
		}
	}

	response, err := s.Client.Do(request)
	if err != nil {
		return nil, err
	}

	if success := s.IsSuccess(response); success {
		if s.SuccessDecoder != nil {
			if err := s.SuccessDecoder(response); err != nil {
				return nil, err
			}
		}
	} else {
		if s.ErrorDecoder != nil {
			if err := s.ErrorDecoder(response); err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}

// Encoder encoders payload into writer. It is only invoked when payload is not-nil.
type Encoder func(writer io.Writer, payload interface{}) error

var (
	// RawEncoder encodes string or bytes into request body.
	RawEncoder Encoder = func(writer io.Writer, payload interface{}) error {
		switch payload.(type) {
		case string:
			if _, err := writer.Write([]byte(payload.(string))); err != nil {
				return err
			}
		case []byte:
			if _, err := writer.Write(payload.([]byte)); err != nil {
				return err
			}
		default:
			return errors.New("invalid type for raw encoder")
		}

		return nil
	}
	// FormEncoder encodes url.Values, map[string][]string, or map[string]string into request body
	// as url encoded form parameters.
	FormEncoder Encoder = func(writer io.Writer, payload interface{}) error {
		switch payload.(type) {
		case url.Values:
			if _, err := writer.Write([]byte(payload.(url.Values).Encode())); err != nil {
				return err
			}
		case map[string][]string:
			if _, err := writer.Write([]byte(url.Values(payload.(map[string][]string)).Encode())); err != nil {
				return err
			}
		case map[string]string:
			u := url.Values{}
			for k, v := range payload.(map[string]string) {
				u.Add(k, v)
			}
			if _, err := writer.Write([]byte(u.Encode())); err != nil {
				return err
			}
		default:
			return errors.New("invalid type for form encoder")
		}
		return nil
	}
	// JSONEncoder encodes payload as JSON into request body.
	JSONEncoder Encoder = func(writer io.Writer, payload interface{}) error {
		return json.NewEncoder(writer).Encode(payload)
	}
	// XMLEncoder encodes payload as XML document into request body.
	XMLEncoder Encoder = func(writer io.Writer, payload interface{}) error {
		return xml.NewEncoder(writer).Encode(payload)
	}
)

// Decoder decodes reader into destination.
type Decoder func(resp *http.Response) error

var (
	// RawDecoder returns a Decoder that decodes the response body as is.
	RawDecoder = func(writer io.Writer) Decoder {
		return func(resp *http.Response) error {
			if writer == nil {
				return errors.New("no writer specified")
			}
			if _, err := io.Copy(writer, resp.Body); err != nil {
				return err
			}
			return nil
		}
	}
	// JSONDecoder returns a Decoder that decodes the response body as JSON into destination.
	JSONDecoder = func(destination interface{}) Decoder {
		return func(resp *http.Response) error {
			return json.NewDecoder(resp.Body).Decode(destination)
		}
	}
	// XMLDecoder returns a Decoder that decodes the response body as XML into destination.
	XMLDecoder = func(destination interface{}) Decoder {
		return func(resp *http.Response) error {
			decoder := xml.NewDecoder(resp.Body)
			decoder.CharsetReader = charset.NewReaderLabel
			return decoder.Decode(destination)
		}
	}
	// FormDecoder returns a Decoder that decodes the response body as application url encoded forms.
	FormDecoder = func(destination *url.Values) Decoder {
		return func(resp *http.Response) error {
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
	// AutoDecoder inspects the Content-Type response header to determine which Decoder to invoke.
	// If content-type starts with "application/json", JSONDecoder is called; If content-type starts
	// with "application/xml", XMLDecoder is called; If content-type starts with "application/x-www-form-urlencoded",
	// FormDecoder is called; Otherwise, RawDecoder is called.
	AutoDecoder = func(destination interface{}) Decoder {
		return func(resp *http.Response) error {
			contentType := resp.Header.Get("Content-Type")
			switch {
			case strings.HasPrefix(contentType, "application/json"):
				return JSONDecoder(destination)(resp)
			case strings.HasPrefix(contentType, "application/xml"):
				return XMLDecoder(destination)(resp)
			case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
				if v, ok := destination.(*url.Values); ok {
					return FormDecoder(v)(resp)
				} else {
					return errors.New("destination not a url.Values")
				}
			default:
				if w, ok := destination.(io.Writer); ok {
					return RawDecoder(w)(resp)
				} else {
					return errors.New("destination not a io.Writer")
				}
			}
		}
	}
)

// IsSuccess returns true if the response is a success.
type IsSuccess func(resp *http.Response) bool

var (
	// DefaultIsSuccess thinks a request has been successful if the response code is 2XX.
	DefaultIsSuccess IsSuccess = func(resp *http.Response) bool {
		return resp.StatusCode < 300 && resp.StatusCode > 199
	}
)
