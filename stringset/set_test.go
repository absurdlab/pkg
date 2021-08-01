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

func testAdd(t *testing.T, p profile) {
	cases := []struct {
		name      string
		setFunc   func() stringset.Interface
		toAdd     string
		sizeAfter int
	}{
		{
			name:      "adding to a nil set has no effect",
			setFunc:   p.nilConstructor,
			toAdd:     "foo",
			sizeAfter: 0,
		},
		{
			name:      "adding to an empty set",
			setFunc:   p.nonNilConstructor,
			toAdd:     "foo",
			sizeAfter: 1,
		},
		{
			name: "adding an existing element",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("foo")
				return set
			},
			toAdd:     "foo",
			sizeAfter: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			set := c.setFunc()
			set.Add(c.toAdd)
			if c.sizeAfter != set.Size() {
				t.Fail()
			}
		})
	}
}

func testRemove(t *testing.T, p profile) {
	cases := []struct {
		name      string
		setFunc   func() stringset.Interface
		toRemove  string
		sizeAfter int
	}{
		{
			name:      "remove from a nil set has no effect",
			setFunc:   p.nilConstructor,
			toRemove:  "foo",
			sizeAfter: 0,
		},
		{
			name:      "remove from an empty set has no effect",
			setFunc:   p.nonNilConstructor,
			toRemove:  "foo",
			sizeAfter: 0,
		},
		{
			name: "remove existing element",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("foo")
				return set
			},
			toRemove:  "foo",
			sizeAfter: 0,
		},
		{
			name: "remove non-existing element",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("foo")
				return set
			},
			toRemove:  "bar",
			sizeAfter: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			set := c.setFunc()
			set.Remove(c.toRemove)
			if c.sizeAfter != set.Size() {
				t.Fail()
			}
		})
	}
}

func testAll(t *testing.T, p profile) {
	var lengthGreaterThanThree stringset.Criteria = func(value string) bool {
		return len(value) > 3
	}

	cases := []struct {
		name     string
		setFunc  func() stringset.Interface
		criteria stringset.Criteria
		expect   bool
	}{
		{
			name:     "nil set fails criteria",
			setFunc:  p.nilConstructor,
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
		{
			name:     "empty set fails criteria",
			setFunc:  p.nonNilConstructor,
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
		{
			name: "matches criteria",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("train")
				return set
			},
			criteria: lengthGreaterThanThree,
			expect:   true,
		},
		{
			name: "does not match criteria",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				stringset.AddAll(set, "foo", "bar")
				return set
			},
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.expect != c.setFunc().All(c.criteria) {
				t.Fail()
			}
		})
	}
}

func testAny(t *testing.T, p profile) {
	var lengthGreaterThanThree stringset.Criteria = func(value string) bool {
		return len(value) > 3
	}

	cases := []struct {
		name     string
		setFunc  func() stringset.Interface
		criteria stringset.Criteria
		expect   bool
	}{
		{
			name:     "nil set fails criteria",
			setFunc:  p.nilConstructor,
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
		{
			name:     "empty set fails criteria",
			setFunc:  p.nonNilConstructor,
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
		{
			name: "matches criteria",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				stringset.AddAll(set, "foo", "train", "bar")
				return set
			},
			criteria: lengthGreaterThanThree,
			expect:   true,
		},
		{
			name: "does not match criteria",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				stringset.AddAll(set, "foo", "bar")
				return set
			},
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.expect != c.setFunc().Any(c.criteria) {
				t.Fail()
			}
		})
	}
}

func testNone(t *testing.T, p profile) {
	var lengthGreaterThanThree stringset.Criteria = func(value string) bool {
		return len(value) > 3
	}

	cases := []struct {
		name     string
		setFunc  func() stringset.Interface
		criteria stringset.Criteria
		expect   bool
	}{
		{
			name:     "nil set fails criteria",
			setFunc:  p.nilConstructor,
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
		{
			name:     "empty set fails criteria",
			setFunc:  p.nonNilConstructor,
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
		{
			name: "matches criteria",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				stringset.AddAll(set, "a", "b", "c")
				return set
			},
			criteria: lengthGreaterThanThree,
			expect:   true,
		},
		{
			name: "does not match criteria",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				stringset.AddAll(set, "a", "train")
				return set
			},
			criteria: lengthGreaterThanThree,
			expect:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.expect != c.setFunc().None(c.criteria) {
				t.Fail()
			}
		})
	}
}

func testArray(t *testing.T, p profile) {
	cases := []struct {
		name      string
		setFunc   func() stringset.Interface
		expectLen int
	}{
		{
			name:      "nil set returns empty array",
			setFunc:   p.nilConstructor,
			expectLen: 0,
		},
		{
			name:      "empty set returns empty array",
			setFunc:   p.nonNilConstructor,
			expectLen: 0,
		},
		{
			name: "non empty set",
			setFunc: func() stringset.Interface {
				set := p.nonNilConstructor()
				set.Add("foo")
				return set
			},
			expectLen: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			set := c.setFunc()
			arr := set.Array()

			if len(arr) != c.expectLen {
				t.Error("unexpected array length")
			}

			if !stringset.ContainsAll(set, arr...) {
				t.Error("set does not contain element")
			}
		})
	}
}

func testOne(t *testing.T, p profile) {
	set := p.nonNilConstructor()
	set.Add("foo")
	if "foo" != set.One() {
		t.Fail()
	}
}
