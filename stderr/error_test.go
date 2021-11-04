package stderr_test

import (
	"errors"
	"fmt"
	"github.com/absurdlab/pkg/stderr"
	"testing"
)

func TestChain_Is(t *testing.T) {
	var (
		foo          = errors.New("foo")
		bar          = errors.New("bar")
		status400    = stderr.Status(400)
		codeNotFound = stderr.Code("not_found")
		messageHello = stderr.Message("hello")
	)

	cases := []struct {
		err    error
		target error
		is     bool
	}{
		{
			err:    stderr.Chain(foo, status400, bar),
			target: stderr.Status(400),
			is:     true,
		},
		{
			err:    stderr.Chain(foo, codeNotFound, bar),
			target: stderr.Code("not_found"),
			is:     true,
		},
		{
			err:    stderr.Chain(foo, codeNotFound, bar),
			target: bar,
			is:     false,
		},
		{
			err:    stderr.Chain(foo, codeNotFound, messageHello),
			target: stderr.Message("hello"),
			is:     false,
		},
		{
			err:    stderr.Chain(foo, codeNotFound, messageHello),
			target: stderr.Params("foo", "bar"),
			is:     false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			expect, actual := c.is, errors.Is(c.err, c.target)
			if expect != actual {
				t.Errorf("want %t, got %t", expect, actual)
			}
		})
	}
}

func TestChain_As(t *testing.T) {
	cases := []struct {
		err error
		run func(t *testing.T, err error)
	}{
		{
			err: stderr.Chain(errors.New("foo"), stderr.Status(400)),
			run: func(t *testing.T, err error) {
				var target *stderr.StatusError
				if !errors.As(err, &target) {
					t.Error("expect to make conversion")
				}

				if expect, actual := 400, target.Status(); expect != actual {
					t.Errorf("expect %d, actual %d", expect, actual)
				}
			},
		},
		{
			err: stderr.Chain(errors.New("foo"), stderr.Code("not_found")),
			run: func(t *testing.T, err error) {
				var target *stderr.CodeError
				if !errors.As(err, &target) {
					t.Error("expect to make conversion")
				}

				if expect, actual := "not_found", target.Code(); expect != actual {
					t.Errorf("expect %s, actual %s", expect, actual)
				}
			},
		},
		{
			err: stderr.Chain(errors.New("foo"), stderr.Params("foo", "bar")),
			run: func(t *testing.T, err error) {
				var target *stderr.ParamsError
				if !errors.As(err, &target) {
					t.Error("expect to make conversion")
				}

				if expect, actual := "bar", target.Params()["foo"]; expect != actual {
					t.Errorf("expect %s, actual %s", expect, actual)
				}
			},
		},
		{
			err: stderr.Chain(errors.New("foo"), stderr.Code("not_found")),
			run: func(t *testing.T, err error) {
				var target *stderr.StatusError
				if errors.As(err, &target) {
					t.Error("expect to not make conversion")
				}
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			c.run(t, c.err)
		})
	}
}
