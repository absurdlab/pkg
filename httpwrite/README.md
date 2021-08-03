# httpwrite

[![.github/workflows/httpwrite.yaml](https://github.com/absurdlab/pkg/actions/workflows/httpwrite.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/httpwrite.yaml)

Render HTTP response in one step.

```shell
go get -u github.com/absurdlab/pkg/httpwrite
```

## Usage

```go
// Use builder to build response, don't worry about http.ResponseWriter call orders.
httpwrite.Render(rw, 
    httpwrite.Options().
    	JSON(payload).
        WithStatus(201).
    	AddHeader("Custom-Header", "hello"),
)
```