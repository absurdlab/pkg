# check

[![.github/workflows/check.yaml](https://github.com/absurdlab/pkg/actions/workflows/check.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/check.yaml)

Fluent validation for common use cases.

```shell
go get -u github.com/absurdlab/pkg/check
```

## Example

```go
// check a string
_ = check.String("david@absurdlab.io").
    Required().
    IsEmail().
    Error()

// check a string slice
_ = check.StringSlice([]string{"foo", "bar", "foo", "bar"}).
    NotEmpty().
    Each(func(elem string) *check.StringCheck {
        return check.String(elem).Required().Enum("foo", "bar")
    }).
    Error(),

// check a URL
_ = check.URLString("https://absurdlab.io").
    HasScheme("https").
    IsNotLocalhost().
    NotHaveQuery().
    NotHaveFragment().
    Error()
```
