# timeplus

Types and utilities to deal with pragmatic use cases of time.

```shell
go get -u github.com/absurdlab/pkg/timeplus
```

## Usage

```go
// Current timestamp
timeplus.Now()

// Arbitrary timestamp
timeplus.On(time.Now().Add(time.Minute))

// Compare timestamp
t1 := timeplus.Now()
t2 := timeplus.Now().AddSecond(timeplus.Second(5))
t1.Before(t2)

// Interop with time.Time
var t = timeplus.Now()
var s = timeplus.Second(5)
t.Time()
s.Duration()

```