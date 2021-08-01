package stringset

func IsEmpty(set Interface) bool {
	return set.Size() == 0
}

func ContainsAll(set Interface, elements ...string) bool {
	for _, each := range elements {
		if !set.Contains(each) {
			return false
		}
	}

	return true
}

func AddAll(set Interface, elements ...string) {
	for _, each := range elements {
		set.Add(each)
	}
}

func RemoveAll(set Interface, elements ...string) {
	for _, each := range elements {
		set.Remove(each)
	}
}

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
