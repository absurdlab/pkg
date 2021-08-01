package stringset

// IsEmpty returns true if this set's size is 0.
func IsEmpty(set Interface) bool {
	return set.Size() == 0
}

// ContainsAll returns true if the set contains all given elements
func ContainsAll(set Interface, elements ...string) bool {
	for _, each := range elements {
		if !set.Contains(each) {
			return false
		}
	}

	return true
}

// Equals returns true s1 and s2 are equal: same size and same elements.
func Equals(s1 Interface, s2 Interface) bool {
	if s1.Size() != s2.Size() {
		return false
	}
	return ContainsAll(s1, s2.Array()...)
}

// AddAll adds all elements into the set. Existing elements are skipped.
func AddAll(set Interface, elements ...string) {
	for _, each := range elements {
		set.Add(each)
	}
}

// RemoveAll removes all elements from the set.
func RemoveAll(set Interface, elements ...string) {
	for _, each := range elements {
		set.Remove(each)
	}
}

// Subset returns true if sub is a subset of the super.
func Subset(super Interface, sub Interface) bool {
	if super.Size() < sub.Size() {
		return false
	}

	return sub.All(func(value string) bool {
		return super.Contains(value)
	})
}

// Coalesce returns the first non-empty set. If all sets are empty, a new Ordered set is returned.
func Coalesce(sets ...Interface) Interface {
	for _, set := range sets {
		if set.Size() > 0 {
			return set
		}
	}
	return NewOrdered()
}
