package check_test

import (
	"github.com/absurdlab/pkg/check"
	"testing"
)

func TestURLChecks(t *testing.T) {
	failures := []struct {
		name string
		err  error
	}{
		{
			name: "not a url",
			err:  check.URLString("\n").Error(),
		},
		{
			name: "not https, but is",
			err:  check.URLString("https://test.org").NotHaveScheme("https").Error(),
		},
		{
			name: "not custom scheme, but is",
			err:  check.URLString("app://deep/link").HasScheme("http", "https").Error(),
		},
		{
			name: "not localhost, but is",
			err:  check.URLString("https://127.0.0.1:8080").IsNotLocalhost().Error(),
		},
		{
			name: "no fragment, but has",
			err:  check.URLString("https://test.org#foobar").NotHaveFragment().Error(),
		},
		{
			name: "no query, but has",
			err:  check.URLString("https://test.org?foo=bar").NotHaveQuery().Error(),
		},
	}

	for _, c := range failures {
		t.Run(c.name, func(t *testing.T) {
			if c.err == nil {
				t.Error("expect test to fail")
			}
		})
	}

	successes := []struct {
		name string
		err  error
	}{
		{
			name: "valid url",
			err:  check.URLString("https://absurdlab.io").Error(),
		},
		{
			name: "meets all criteria",
			err: check.URLString("https://absurdlab.io").
				HasScheme("https").
				IsNotLocalhost().
				NotHaveQuery().
				NotHaveFragment().
				Error(),
		},
	}

	for _, c := range successes {
		t.Run(c.name, func(t *testing.T) {
			if c.err != nil {
				t.Error(c.err)
			}
		})
	}
}
