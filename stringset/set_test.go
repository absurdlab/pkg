package stringset_test

import (
	"github.com/absurdlab/pkg/stringset"
	"testing"
)

type profile struct {
	nonNilConstructor func() stringset.Interface
	nilConstructor    func() stringset.Interface
}

func testSize(t *testing.T, p profile) {
	cases := []struct {
		name        string
		setFunc     func() stringset.Interface
		constructor func() stringset.Interface
		mod         func(set stringset.Interface)
		assert      func(t *testing.T, size int)
	}{
		{
			name:    "nil set should have 0 size",
			setFunc: p.nilConstructor,
			assert: func(t *testing.T, size int) {
				if size != 0 {
					t.Fail()
				}
			},
		},
		{
			name:    "empty set should have 0 size",
			setFunc: p.nonNilConstructor,
			assert: func(t *testing.T, size int) {
				if size != 0 {
					t.Fail()
				}
			},
		},
		{
			name: "set with one element should have 1 size",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("foo")
				return set
			},
			assert: func(t *testing.T, size int) {
				if size != 1 {
					t.Fail()
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.assert(t, c.setFunc().Size())
		})
	}
}

func testContains(t *testing.T, p profile) {
	cases := []struct {
		name     string
		setFunc  func() stringset.Interface
		element  string
		contains bool
	}{
		{
			name:     "nil set never contains anything",
			setFunc:  p.nilConstructor,
			element:  "foo",
			contains: false,
		},
		{
			name:     "empty set does not contain foo",
			setFunc:  p.nonNilConstructor,
			element:  "foo",
			contains: false,
		},
		{
			name: "set with foo contains foo",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("foo")
				return set
			},
			element:  "foo",
			contains: true,
		},
		{
			name: "set with bar does not contain foo",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("bar")
				return set
			},
			element:  "foo",
			contains: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.contains != c.setFunc().Contains(c.element) {
				t.Fail()
			}
		})
	}
}
