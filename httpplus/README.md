# httpplus

Light utilities to make it easier to work with http package.

```shell
go get -u github.com/absurdlab/pkg/httpplus
```

## Make Request

```go
// Making a http request and receive decoded response. 
var (
    payload = &Payload{Message: "hello world"}
    successReply = new(SuccessReply)
    errorReply = new(ErrorReply)
)

spec := &httputil.RequestSpec{
    Method: http.MethodPost,
    URL: "https://httpbin.org/post",
    Headers: map[string]string{"Content-Type", "application/json"},
    Payload: payload,
    Encoder: httpplus.JSONEncoder,
    SuccessDecoder: httpplus.JSONDecoder(successReply),
    ErrorDecoder: httpplus.JSONDecoder(errorReply),
    IsSuccess: httpplus.DefaultIsSuccess,
}

// successReply, or errorReply should have been populated
httpResponse, err = httpplus.MakeRequest(ctx, spec)
```

### Encoder

`Encoder` is responsible for encoding request payload into desired format, included encoders are:
- `RawEncoder`: encodes `[]byte` or `string` as is.
- `JSONEncoder`: encodes structure as JSON.
- `XMLEncoder`: encodes structure as XML.
- `FormEncoder`: encodes `url.Values`, `map[string][]string`, or `map[string]string` as URL encoded form.

**Note that encoders will not automatically include `Content-Type` headers.**

### Decoder

`Decoder` is responsible for decoding response payload into desired structure, included decoders are:
- `RawDecoder`: decodes into `io.Writer` (i.e. `bytes.Buffer`, `strings.Builder`)
- `JSONDecoder`: decodes into the given structure as JSON.
- `XMLDecoder`: decodes into the given structure as XML.
- `FormDecoder`: decodes into `url.Values` as url encoded parameters.
- `AutoDecoder`: call the above encoders depending on `Content-Type`, and fallback to `RawDecoder`.

### IsSuccess

`IsSuccess` is a criteria to decide whether the operation was successful. Successful response will invoke
`SuccessDecoder` if defined; Error response will invoke `ErrorDecoder` if defined. 

By default, all `2XX` status codes are deemed as success, everything else is an error.