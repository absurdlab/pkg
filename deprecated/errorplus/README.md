# errorplus

[![.github/workflows/errorplus.yaml](https://github.com/absurdlab/pkg/actions/workflows/errorplus.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/errorplus.yaml)

Standard mechanism to create a chain of error and render final error response.

```shell
go get -u github.com/absurdlab/pkg/errorplus
```

## Example Usage

```go
// At the repository level
userNotFoundErr := &UserNotFoundError{ID: "f02fb2003ea04d5"}
return errorplus.Decorate(userNotFoundErr).
    Wrap(sql.ErrNoRows).
    StatusHint(http.StatusNotFound)

// Later, at some service
return errorplus.Wrap(err).Field("operation", "find_user")

// Later, at another service. We have more context.
userNotFoundErr := new(UserNotFoundError)
_ = errors.As(err, &userNotFoundErr)
return errorplus.Wrap(err).
    DecorateF("invalid_group/user_not_found").
    Field("operation", "update_group").
    Field("group_id", "2c515ff8c845").
    Field("user_id", userNotFoundErr.ID)
    StatusHint(http.StatusBadRequest)

// finally, at the handler level
details := errorplus.GetDetailStack(err)
json.NewEncoder(rw).Encode(errorplus.APIError{
    Status: errorplus.GetStatusHint(err),
    Timestamp: timeplus.Now().Ref(),
    RequestID: requestIdFromContext(ctx),
    TraceID: traceIdFromContext(ctx),
    Error: errorplus.GetMessage(err),
    Detail: details[0],
    Stack: details,
})
```