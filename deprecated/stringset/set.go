package stringset

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

// Interface is the main abstraction for a string set data structure.
type Interface interface {
	json.Marshaler
	json.Unmarshaler
	yaml.Marshaler
	yaml.Unmarshaler
	yaml.IsZeroer
	sql.Scanner
	driver.Valuer

	// Size returns the size of the set. If set is nil, returns 0.
	Size() int
	// Contains returns true if the queried element is in the set.
	Contains(element string) bool
	// Add adds the given element to the set, if it does not already exist. Adding to a nil set has no effect.
	Add(element string)
	// Remove removes the given element from the set, if it exists.
	Remove(element string)
	// All returns true if all elements in the set matches the given criteria. Empty set always returns false.
	All(criteria Criteria) bool
	// Any returns true if at least one element in the set matches given criteria. Empty set always returns false.
	Any(criteria Criteria) bool
	// None returns true if none elements in the set matches given criteria. Empty set always returns false.
	None(criteria Criteria) bool
	// Array returns a copied slice of the elements in this set in ordered fashion. Subsequent operations on the returned
	// slice will not affect the set in any way. If set is nil or empty, an empty slice is returned.
	Array() []string
	// One returns one element from this set. It is usually invoked when Size is 1, but can be implemented for
	// any Size greater than 0. If set is empty, implementations should panic.
	One() string
}

// Criteria is a function that help selects an element in the set.
type Criteria func(value string) bool
