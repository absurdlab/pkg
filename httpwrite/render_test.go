package httpwrite_test

import (
	"github.com/absurdlab/pkg/httpwrite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRender(t *testing.T) {
	cases := []struct {
		name    string
		options *httpwrite.WriteOptions
		assert  func(t *testing.T, rr *httptest.ResponseRecorder, err error)
	}{
		{
			name:    "status only",
			options: httpwrite.Options().WithStatus(http.StatusCreated),
			assert: func(t *testing.T, rr *httptest.ResponseRecorder, err error) {
				if err != nil {
					t.Error(err)
				}

				if rr.Code != http.StatusCreated {
					t.Error("expect status to be 201")
				}
			},
		},
		{
			name:    "plain text",
			options: httpwrite.Options().PlainText("hello"),
			assert: func(t *testing.T, rr *httptest.ResponseRecorder, err error) {
				if err != nil {
					t.Error(err)
				}

				if rr.Code != http.StatusOK {
					t.Error("expect status to be 200")
				}

				if rr.Header().Get("Content-Type") != "text/plain" {
					t.Error("expect content type to be text/plain")
				}

				if rr.Body.String() != "hello" {
					t.Error("expect body to be hello")
				}
			},
		},
		{
			name:    "json",
			options: httpwrite.Options().JSON(map[string]string{"message": "hello"}).WithStatus(http.StatusCreated),
			assert: func(t *testing.T, rr *httptest.ResponseRecorder, err error) {
				if err != nil {
					t.Error(err)
				}

				if rr.Code != http.StatusCreated {
					t.Error("expect status to be 201")
				}

				if rr.Header().Get("Content-Type") != "application/json" {
					t.Error("expect content type to be application/json")
				}

				if rr.Body.String() != "{\"message\":\"hello\"}\n" {
					t.Error("expect body to be json")
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := httpwrite.Render(rr, c.options)
			c.assert(t, rr, err)
		})
	}
}
