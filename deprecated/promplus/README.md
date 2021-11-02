# promplus

[![.github/workflows/promplus.yaml](https://github.com/absurdlab/pkg/actions/workflows/promplus.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/promplus.yaml)

Out-of-box HTTP middlewares to collect common prometheus metrics

```shell
go get -u github.com/absurdlab/pkg/promplus
```

## Metrics

- :watch: handler duration
- :ledger: handler visit count
- :ok: handler response status count

## Usage

```go
target := &promplus.Target{
    Namespace: "my_company",
    Service: "my_app",
}

var h = myHandler()
h = target.Duration()(h)
h = target.VisitCount()(h)
h = target.StatusCount()(h)
```