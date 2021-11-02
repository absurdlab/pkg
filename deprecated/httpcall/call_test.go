package httpcall_test

import (
	"context"
	"encoding/xml"
	"github.com/absurdlab/pkg/httpcall"
	"net/http"
	"strings"
	"testing"
)

func TestMake(t *testing.T) {
	cases := []struct {
		name        string
		prepareFunc func() (*body, *httpcall.CallOptions)
		assert      func(t *testing.T, resp *http.Response, r *body, err error)
	}{
		{
			name: "get raw response",
			prepareFunc: func() (*body, *httpcall.CallOptions) {
				var sb strings.Builder
				return &body{success: &sb}, httpcall.Options().
					WithURL("https://httpbin.org/encoding/utf8").
					AddHeaders("Accept", "text/html").
					ToPlainSuccess(&sb)
			},

			assert: func(t *testing.T, resp *http.Response, r *body, err error) {
				if err != nil {
					t.Error(err)
				}

				body := r.success.(*strings.Builder).String()
				if len(body) == 0 {
					t.Error("expected non-empty payload from robot.txt")
				}
			},
		},
		{
			name: "get json response",
			prepareFunc: func() (*body, *httpcall.CallOptions) {
				var reply httpBinJSONResponse
				return &body{success: &reply}, httpcall.Options().
					WithURL("https://httpbin.org/get").
					ToJSONSuccess(&reply)
			},
			assert: func(t *testing.T, resp *http.Response, r *body, err error) {
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
			name: "get xml response",
			prepareFunc: func() (*body, *httpcall.CallOptions) {
				var reply httpBinXMLResponse
				return &body{success: &reply}, httpcall.Options().
					WithURL("https://httpbin.org/xml").
					AddHeaders("Accept", "application/xml").
					ToXMLSuccess(&reply)
			},
			assert: func(t *testing.T, resp *http.Response, r *body, err error) {
				if err != nil {
					t.Error(err)
				}

				if len(r.success.(*httpBinXMLResponse).Title) == 0 {
					t.Error("expect slideshow title in xml")
				}
			},
		},
		{
			name: "post json message",
			prepareFunc: func() (*body, *httpcall.CallOptions) {
				var reply httpBinJSONResponse
				return &body{success: &reply}, httpcall.Options().
					POST("https://httpbin.org/post").
					JSON(map[string]string{"message": "hello"}).
					ToJSONSuccess(&reply)
			},
			assert: func(t *testing.T, resp *http.Response, r *body, err error) {
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
			name: "post form",
			prepareFunc: func() (*body, *httpcall.CallOptions) {
				var reply httpBinJSONResponse
				return &body{success: &reply}, httpcall.Options().
					POST("https://httpbin.org/anything").
					Form(map[string]string{"foo": "bar"}).
					ToJSONSuccess(&reply)
			},
			assert: func(t *testing.T, resp *http.Response, r *body, err error) {
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
			prepareFunc: func() (*body, *httpcall.CallOptions) {
				var sb strings.Builder
				return &body{success: &sb}, httpcall.Options().
					WithURL("https://httpbin.org/status/201").
					AddHeaders("Accept", "text/plain").
					IsSuccessWhenStatus(http.StatusCreated).
					ToSuccess(func(resp *http.Response) error {
						(&sb).WriteString("success")
						return nil
					})
			},
			assert: func(t *testing.T, resp *http.Response, r *body, err error) {
				if err != nil {
					t.Error(err)
				}

				body := r.success.(*strings.Builder).String()
				if body != "success" {
					t.Error("expected success decoder to be called")
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bodies, options := c.prepareFunc()
			response, err := httpcall.Make(context.Background(), options)
			c.assert(t, response, bodies, err)
		})
	}
}

type body struct {
	success interface{}
	error   interface{}
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
