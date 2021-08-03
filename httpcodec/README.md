# httpcodec

easy HTTP encoding and decoding.

```shell
go get -u github.com/absurdlab/pkg/httpcodec
```

## Features

- Raw/Plain codec
- JSON codec
- XML codec
- Form codec

## Usage

```go
// encode
httpcodec.EncodeRaw(rw, "hello")
httpcodec.EncodeJSON(rw, jsonPayload{Message: "hello"})
httpcodec.EncodeXML(rw, xmlPayload{Message: "hello"})
httpcodec.EncodeForm(rw, map[string]string{"greeting": "hello"})

// decode JSON
httpcodec.DecodeRaw(&stringBuilder)(httpResp)
httpcodec.DecodeJSON(payload)(httpResp)
httpcodec.DecodeXML(payload)(httpResp)
httpcodec.DecodeForm(&urlValues)(httpResp)
```