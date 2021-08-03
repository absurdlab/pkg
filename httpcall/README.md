# httpcall

[![.github/workflows/httpcall.yaml](https://github.com/absurdlab/pkg/actions/workflows/httpcall.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/httpcall.yaml)

Calling HTTP endpoints is a bit easier now.

```shell
go get -u github.com/absurdlab/pkg/httpcall
```

## Usage

```go
var (
    payload = Greeting{Message: "hello"}
    success SuccessReply
    failure ErrorReply
    
    options = httpcall.Options().
    	POST("https://httpbin.org/anything").
    	JSON(payload).
        ToJSONSuccess(&success).
        ToJSONError(&failure)
)

_, _ = httpcall.Make(context.Background(), options)
```