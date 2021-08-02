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

// MustURLWithQuery encodes the query parameters to the query component of the url string. The provided url string
// must be valid, otherwise the function panics. If the url string already contains some query parameters, the added
// parameters will add, instead of overwrite.
func MustURLWithQuery(urlStr string, queryParams url.Values) string {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		panic(err)
	}

	p := u.Query()
	for k, vs := range queryParams {
		for _, v := range vs {
			p.Add(k, v)
		}
	}

	u.RawQuery = p.Encode()

	return u.String()
}
