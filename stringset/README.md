# stringset

[![.github/workflows/stringset.yaml](https://github.com/absurdlab/pkg/actions/workflows/stringset.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/stringset.yaml)

String based set data structure. 

```shell
go get -u github.com/absurdlab/pkg/stringset
```

## Features

- Ordered and non-Ordered implementations
- Utility functions
- JSON marshalling ready
- YAML marshalling ready
- SQL marshalling ready

## Usage

```go
// create a new ordered set with values
stringset.NewOrderedWith("foo", "bar")

// create a new ordered set from space delimited values
stringset.NewOrderedBySpace("foo bar")
```