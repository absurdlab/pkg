package httpplus_test

import (
	"fmt"
	"github.com/absurdlab/pkg/httpplus"
	"net/url"
	"testing"
)

func TestURLValues(t *testing.T) {
	v := httpplus.URLValues(
		"Content-Type", "application/json",
		"Accept", "application/json",
	)

	if len(v) != 2 {
		t.Error("expect url.Values to have 2 entries")
	}

	if contentType := v.Get("Content-Type"); len(contentType) == 0 {
		t.Error("expect url.Values to have Content-Type")
	}

	if accept := v.Get("Accept"); len(accept) == 0 {
		t.Error("expect url.Values to have Accept")
	}
}

func TestMustURLWithQuery(t *testing.T) {
	cases := []struct {
		urlStr string
		params url.Values
		expect string
	}{
		{
			urlStr: "https://test.org",
			params: httpplus.URLValues("foo", "bar"),
			expect: "https://test.org?foo=bar",
		},
		{
			urlStr: "https://test.org?foo=bar",
			params: httpplus.URLValues("foo", "bar"),
			expect: "https://test.org?foo=bar&foo=bar",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("TestMustURLWithQuery/%d", i), func(t *testing.T) {
			actual := httpplus.MustURLWithQuery(c.urlStr, c.params)
			if c.expect != actual {
				t.Error("expect encoded url to match")
			}
		})
	}
}
