package httpwrite

import "net/http"

// Render renders the response according to WriteOptions.
func Render(rw http.ResponseWriter, options *WriteOptions) error {
	options.sanitize() // sanitize just in case

	if len(options.headers) > 0 {
		for k, v := range options.headers {
			rw.Header().Add(k, v)
		}
	}

	rw.WriteHeader(options.status)

	if options.body != nil && options.encoder != nil {
		if err := options.encoder(rw, options.body); err != nil {
			return err
		}
	}

	return nil
}
