package httpplus_test

import (
	"context"
	"encoding/xml"
	"github.com/absurdlab/pkg/httpplus"
	"net/http"
	"strings"
	"testing"
)

func TestRequestSpec_Execute(t *testing.T) {
	cases := []struct {
		name      string
		responses *responses
		specFunc  func(r *responses) *httpplus.RequestSpec
		assert    func(t *testing.T, r *responses, err error)
	}{
		{
			name: "get raw response",
			responses: &responses{
				success: new(strings.Builder),
			},
			specFunc: func(r *responses) *httpplus.RequestSpec {
				return &httpplus.RequestSpec{
					URL:            "https://httpbin.org/encoding/utf8",
					Headers:        map[string]string{"Accept": "text/html"},
					SuccessDecoder: httpplus.RawDecoder(r.success.(*strings.Builder)),
				}
			},
			assert: func(t *testing.T, r *responses, err error) {
				if err != nil {
					t.Error(err)
				}

				body := r.success.(*strings.Builder).String()
				if len(body) == 0 {
					t.Error("expected non-empty payload from robot.txt")
				}

				t.Log(body)
			},
		},
		{
			name:      "get json response",
			responses: &responses{success: new(httpBinJSONResponse)},
			specFunc: func(r *responses) *httpplus.RequestSpec {
				return &httpplus.RequestSpec{
					URL:            "https://httpbin.org/get",
					SuccessDecoder: httpplus.JSONDecoder(r.success),
				}
			},
			assert: func(t *testing.T, r *responses, err error) {
				if err != nil {
					t.Error(err)
				}

				body := r.success.(*httpBinJSONResponse)
				if body.URL != "https://httpbin.org/get" {
					t.Error("expected payload body URL to be 'https://httpbin.org/get'")
				}
			},
		},
		{
			name:      "get xml response",
			responses: &responses{success: new(httpBinXMLResponse)},
			specFunc: func(r *responses) *httpplus.RequestSpec {
				return &httpplus.RequestSpec{
					URL:            "https://httpbin.org/xml",
					Headers:        map[string]string{"Accept": "application/xml"},
					SuccessDecoder: httpplus.XMLDecoder(r.success),
				}
			},
			assert: func(t *testing.T, r *responses, err error) {
				if err != nil {
					t.Error(err)
				}

				if len(r.success.(*httpBinXMLResponse).Title) == 0 {
					t.Error("expect slideshow title in xml")
				}
			},
		},
		{
			name:      "post json message",
			responses: &responses{success: new(httpBinJSONResponse)},
			specFunc: func(r *responses) *httpplus.RequestSpec {
				return &httpplus.RequestSpec{
					Method: http.MethodPost,
					URL:    "https://httpbin.org/post",
					Headers: map[string]string{
						"Accept":       "application/json",
						"Content-Type": "application/json",
					},
					Payload:        map[string]string{"message": "hello"},
					Encoder:        httpplus.JSONEncoder,
					SuccessDecoder: httpplus.JSONDecoder(r.success),
				}
			},
			assert: func(t *testing.T, r *responses, err error) {
				if err != nil {
					t.Error(err)
				}

				body := r.success.(*httpBinJSONResponse)
				if body.JSON["message"] != "hello" {
					t.Error("expected http bin json post to echo the same body")
				}
			},
		},
		{
			name:      "post form",
			responses: &responses{success: new(httpBinJSONResponse)},
			specFunc: func(r *responses) *httpplus.RequestSpec {
				return &httpplus.RequestSpec{
					Method: http.MethodPost,
					URL:    "https://httpbin.org/anything",
					Headers: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					Payload:        map[string]string{"foo": "bar"},
					Encoder:        httpplus.FormEncoder,
					SuccessDecoder: httpplus.JSONDecoder(r.success),
				}
			},
			assert: func(t *testing.T, r *responses, err error) {
				if err != nil {
					t.Error(err)
				}

				if r.success.(*httpBinJSONResponse).Form["foo"] != "bar" {
					t.Error("expect httpbin to respond with form values")
				}
			},
		},
		{
			name: "custom success condition",
			responses: &responses{
				success: new(strings.Builder),
			},
			specFunc: func(r *responses) *httpplus.RequestSpec {
				return &httpplus.RequestSpec{
					URL: "https://httpbin.org/status/201",
					Headers: map[string]string{
						"Accept": "text/plain",
					},
					SuccessDecoder: func(_ *http.Response) error {
						r.success.(*strings.Builder).WriteString("success")
						return nil
					},
					IsSuccess: func(resp *http.Response) bool {
						return resp.StatusCode == http.StatusCreated
					},
				}
			},
			assert: func(t *testing.T, r *responses, err error) {
				if err != nil {
					t.Error(err)
				}

				body := r.success.(*strings.Builder).String()
				if body != "success" {
					t.Error("expected success decoder to be called")
				}
			},
		},
		{
			name:      "auto decoder",
			responses: &responses{success: new(httpBinJSONResponse)},
			specFunc: func(r *responses) *httpplus.RequestSpec {
				return &httpplus.RequestSpec{
					URL:            "https://httpbin.org/get",
					SuccessDecoder: httpplus.AutoDecoder(r.success),
				}
			},
			assert: func(t *testing.T, r *responses, err error) {
				if err != nil {
					t.Error(err)
				}

				body := r.success.(*httpBinJSONResponse)
				if body.URL != "https://httpbin.org/get" {
					t.Error("expected payload body URL to be 'https://httpbin.org/get'")
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var err error
			spec := c.specFunc(c.responses)
			c.responses.raw, err = httpplus.MakeRequest(context.Background(), spec)
			c.assert(t, c.responses, err)
		})
	}
}

type responses struct {
	success interface{}
	error   interface{}
	raw     *http.Response
}

type httpBinJSONResponse struct {
	Args    map[string]string      `json:"args"`
	Form    map[string]string      `json:"form"`
	Headers map[string]string      `json:"headers"`
	Origin  string                 `json:"origin"`
	JSON    map[string]interface{} `json:"json"`
	URL     string                 `json:"url"`
}

type httpBinXMLResponse struct {
	XMLName xml.Name `xml:"slideshow"`
	Title   string   `xml:"title,attr"`
}
