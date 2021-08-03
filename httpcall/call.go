package httpcall

import (
	"context"
	"net/http"
)

// Make makes the http call according to the call context options and returns the http.Response.
func Make(ctx context.Context, context *callContext) (*http.Response, error) {
	context.sanitize() // sanitize just in case.

	url, err := context.urlString()
	if err != nil {
		return nil, err
	}

	body, err := context.body()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, context.method, url, body)
	if err != nil {
		return nil, err
	}

	response, err := context.client.Do(request)
	if err != nil {
		return nil, err
	}

	success := context.isSuccess(response)
	if success {
		if context.successDecoder != nil {
			err = context.successDecoder(response)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if context.errorDecoder != nil {
			err = context.errorDecoder(response)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}
