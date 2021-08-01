package stringset_test

import (
	"github.com/absurdlab/pkg/stringset"
	"testing"
)

var unorderedSetTestProfile = profile{
	nonNilConstructor: func() stringset.Interface {
		return stringset.NewUnOrdered()
	},
	nilConstructor: func() stringset.Interface {
		var set *stringset.UnOrdered
		return set
	},
}

func TestUnOrdered_Size(t *testing.T) {
	testSize(t, unorderedSetTestProfile)
}

func TestUnOrdered_Contains(t *testing.T) {
	testContains(t, unorderedSetTestProfile)
}

func TestUnOrdered_Add(t *testing.T) {
	testAdd(t, unorderedSetTestProfile)
}

func TestUnOrdered_Remove(t *testing.T) {
	testRemove(t, unorderedSetTestProfile)
}

func TestUnOrdered_All(t *testing.T) {
	testAll(t, unorderedSetTestProfile)
}

func TestUnOrdered_Any(t *testing.T) {
	testAny(t, unorderedSetTestProfile)
}

func TestUnOrdered_None(t *testing.T) {
	testNone(t, unorderedSetTestProfile)
}

func TestUnOrdered_Array(t *testing.T) {
	testArray(t, unorderedSetTestProfile)
}

func TestUnOrdered_One(t *testing.T) {
	testOne(t, unorderedSetTestProfile)
}
