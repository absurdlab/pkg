package check_test

import (
	"github.com/absurdlab/pkg/check"
	"testing"
)

func TestStringSliceChecks(t *testing.T) {
	failures := []struct {
		name string
		err  error
	}{
		{
			name: "not nil but is",
			err:  check.StringSlice(nil).NotNil().Error(),
		},
		{
			name: "not empty but is",
			err:  check.StringSlice([]string{}).NotEmpty().Error(),
		},
		{
			name: "unsatisfied element",
			err: check.StringSlice([]string{"foo", "bar", "baz"}).
				Each(func(elem string) *check.StringCheck {
					return check.String(elem).Required().Enum("foo", "bar")
				}).
				Error(),
		},
		{
			name: "any",
			err: check.StringSlice([]string{"foo", "bar", "baz"}).
				Each(func(elem string) *check.StringCheck {
					return check.String(elem).Equals("abc")
				}).
				Error(),
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
			name: "not empty",
			err:  check.StringSlice([]string{"foo"}).NotEmpty().Error(),
		},
		{
			name: "optional",
			err:  check.StringSlice(nil).Optional().Error(),
		},
		{
			name: "satisfies element",
			err: check.StringSlice([]string{"foo", "bar", "foo", "bar"}).
				Each(func(elem string) *check.StringCheck {
					return check.String(elem).Required().Enum("foo", "bar")
				}).
				Error(),
		},
		{
			name: "any",
			err: check.StringSlice([]string{"foo", "bar", "baz"}).
				Each(func(elem string) *check.StringCheck {
					return check.String(elem).Equals("foo")
				}).
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
