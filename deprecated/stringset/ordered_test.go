package stringset_test

import (
	"github.com/absurdlab/pkg/stringset"
	"testing"
)

var orderedSetTestProfile = profile{
	nonNilConstructor: func() stringset.Interface {
		return stringset.NewOrdered()
	},
	nilConstructor: func() stringset.Interface {
		var set *stringset.Ordered
		return set
	},
}

func TestOrdered_Size(t *testing.T) {
	testSize(t, orderedSetTestProfile)
}

func TestOrdered_Contains(t *testing.T) {
	testContains(t, orderedSetTestProfile)
}

func TestOrdered_Add(t *testing.T) {
	testAdd(t, orderedSetTestProfile)
}

func TestOrdered_Remove(t *testing.T) {
	testRemove(t, orderedSetTestProfile)
}

func TestOrdered_All(t *testing.T) {
	testAll(t, orderedSetTestProfile)
}

func TestOrdered_Any(t *testing.T) {
	testAny(t, orderedSetTestProfile)
}

func TestOrdered_None(t *testing.T) {
	testNone(t, orderedSetTestProfile)
}

func TestOrdered_Array(t *testing.T) {
	testArray(t, orderedSetTestProfile)
}

func TestOrdered_One(t *testing.T) {
	testOne(t, orderedSetTestProfile)
}

func TestOrdered_MarshalJSON(t *testing.T) {
	cases := []struct {
		name    string
		setFunc func() stringset.Interface
		expect  string
	}{
		{
			name:    "nil set",
			setFunc: orderedSetTestProfile.nilConstructor,
			expect:  "null",
		},
		{
			name:    "empty set",
			setFunc: orderedSetTestProfile.nonNilConstructor,
			expect:  "[]",
		},
		{
			name: "non-empty set",
			setFunc: func() stringset.Interface {
				return stringset.NewOrderedWith("foo", "bar")
			},
			expect: `["foo","bar"]`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, err := c.setFunc().MarshalJSON()
			if err != nil {
				t.Error(err)
			}

			if c.expect != string(raw) {
				t.Error("json mismatch")
			}
		})
	}
}

func TestOrdered_MarshalYAML(t *testing.T) {
	cases := []struct {
		name    string
		setFunc func() stringset.Interface
		expect  string
	}{
		{
			name:    "nil set",
			setFunc: orderedSetTestProfile.nilConstructor,
			expect:  "null\n",
		},
		{
			name:    "empty set",
			setFunc: orderedSetTestProfile.nonNilConstructor,
			expect:  "[]\n",
		},
		{
			name: "non-empty set",
			setFunc: func() stringset.Interface {
				return stringset.NewOrderedWith("foo", "bar")
			},
			expect: "- foo\n- bar\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, err := c.setFunc().MarshalYAML()
			if err != nil {
				t.Error(err)
			}
			if c.expect != string(raw.([]byte)) {
				t.Fail()
			}
		})
	}
}

func TestNewOrderedBySpace(t *testing.T) {
	if ok := stringset.Equals(
		stringset.NewOrderedBySpace("foo bar"),
		stringset.NewOrderedWith("foo", "bar"),
	); !ok {
		t.Fail()
	}
}
