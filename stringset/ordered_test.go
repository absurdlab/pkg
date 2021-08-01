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
