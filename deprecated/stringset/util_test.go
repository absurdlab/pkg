package stringset_test

import (
	"github.com/absurdlab/pkg/stringset"
	"testing"
)

func TestCoalesce(t *testing.T) {
	var s0 stringset.Interface = stringset.NewOrdered()
	var s1 stringset.Interface = stringset.NewOrderedWith("foo")
	if !stringset.Equals(stringset.Coalesce(s0, s1), s1) {
		t.Fail()
	}
}

func TestSubset(t *testing.T) {
	var s0 stringset.Interface = stringset.NewOrderedWith("foo", "bar")
	var s1 stringset.Interface = stringset.NewOrderedWith("foo", "bar", "foobar")
	var s2 stringset.Interface = stringset.NewOrderedWith("foo", "baz")

	if ok := stringset.Subset(s1, s0); !ok {
		t.Error("s0 should be subset of s1")
	}

	if ok := stringset.Subset(s2, s0); ok {
		t.Error("s2 should not be subset of s0")
	}
}
