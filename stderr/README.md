# stderr

Yet another attempt to enrich Go's error UX. Chain pieces of error context together as we go up the callstack, and easily render.
Receive and reconstruct the entire error context of the caller at the other end of the wire.

## Install

```bash
go get -u github.com/absurdlab/pkg/stderr
```

## 101

There are four type of errors:
- `*stderr.StatusError`: carries http status value
- `*stderr.CodeError`: carries an error code
- `*stderr.MessageError`: carries a human readable error message
- `*stderr.ParamsError`: carries key value pairs of error context
- `*stderr.GenericError`: wraps a generic error

Errors can be chained together using `Chain`:

```go
// earlier in callstack
err := stderr.Chain(stderr.Message("something went wrong"), errors.New("unexpected error"))

// later in callstack, when we have a bit more context
err := stderr.Chain(stderr.CodeError("invalid_state"), stderr.Params("state", myState), err)

// finally when rendering API response
err := stderr.Chain(stderr.Status(500), err)
```

We can convert the error to a `View` that is ready to render:

```go
view := stderr.ToView(err)
response.JSON(view.Status, view)
```

On the receiving end, we can parse the view and restore the entire error chain:

```go
err := stderr.FromView(view)

var status *stderr.StatusError
errors.As(err, &status)

var params *stderr.ParamsError
errors.As(err, &params)
```