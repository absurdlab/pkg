package jwa

const None = "none"

// IsNoneOrEmpty returns true if the given algorithm is None or is empty.
func IsNoneOrEmpty(alg string) bool {
	return len(alg) == 0 || alg == None
}

// IsDefined is the opposite of IsNoneOrEmpty.
func IsDefined(alg string) bool {
	return !IsNoneOrEmpty(alg)
}
