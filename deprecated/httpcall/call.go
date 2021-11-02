package httpcall

import (
	"context"
	"net/http"
)

// Make makes the http call according to the call context options and returns the http.Response.
func Make(ctx context.Context, options *CallOptions) (*http.Response, error) {
	options.sanitize() // sanitize just in case.

	url, err := options.urlString()
	if err != nil {
		return nil, err
	}

	body, err := options.body()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, options.method, url, body)
	if err != nil {
		return nil, err
	}

	if len(options.headers) > 0 {
		for k, v := range options.headers {
			request.Header.Add(k, v)
		}
	}

	response, err := options.client.Do(request)
	if err != nil {
		return nil, err
	}

	success := options.isSuccess(response)
	if success {
		if options.successDecoder != nil {
			err = options.successDecoder(response)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if options.errorDecoder != nil {
			err = options.errorDecoder(response)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}
