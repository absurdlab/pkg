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
