package check_test

import (
	"github.com/absurdlab/pkg/check"
	"net/url"
	"testing"
)

func TestStringChecks(t *testing.T) {
	failures := []struct {
		name string
		err  error
	}{
		{
			name: "required but absent",
			err:  check.String("").Required().Error(),
		},
		{
			name: "not within enum",
			err:  check.String("foo").Enum("bar", "baz").Error(),
		},
		{
			name: "not equal",
			err:  check.String("foo").Equals("bar").Error(),
		},
		{
			name: "not an email",
			err:  check.String("foo").IsEmail().Error(),
		},
		{
			name: "not valid json",
			err:  check.String("foo").IsValidJSON().Error(),
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
			name: "required and present",
			err:  check.String("foo").Required().Error(),
		},
		{
			name: "optional and absent",
			err:  check.String("").Optional().Error(),
		},
		{
			name: "equals",
			err:  check.String("foo").Equals("foo").Error(),
		},
		{
			name: "equal case insensitive",
			err:  check.String("foo").EqualsCaseInsensitive("FOO").Error(),
		},
		{
			name: "enum",
			err:  check.String("foo").Enum("foo", "bar").Error(),
		},
		{
			name: "email",
			err:  check.String("david@absurdlab.io").IsEmail().Error(),
		},
		{
			name: "url",
			err: check.String("https://absurdlab.io").
				IsURL(func(value *url.URL) *check.URLCheck {
					return check.URL(value).
						HasScheme("https").
						IsNotLocalhost().
						NotHaveQuery().
						NotHaveFragment()
				}).
				Error(),
		},
		{
			name: "valid json",
			err:  check.String(`{"foo": "bar"}`).IsValidJSON().Error(),
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
