package stringset

// Criteria is a function that help selects an element in the set.
type Criteria func(value string) bool
