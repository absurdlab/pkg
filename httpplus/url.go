package httpplus

import "net/url"

// URLValues creates url.Values using key value pairs.
func URLValues(kvs ...string) url.Values {
	if len(kvs)%2 != 0 {
		panic("kvs must be supplied in pairs")
	}

	values := url.Values{}
	for i := 0; i < len(kvs); i += 2 {
		values.Add(kvs[i], kvs[i+1])
	}

	return values
}
