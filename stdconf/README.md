# stdconf

Opinionated way to load configuration from various sources.

## Install

```bash
go get -u github.com/absurdlab/pkg/stdconf
```

## Usage

```go
v, _ := stdconf.Parse(
    stdconf.WithNewFunc(newConfig),
    stdconf.WithMergoOptions(mergo.WithOverride),
    stdconf.WithSources(
        stdconf.FromYAMLFile("config.yaml"),
        stdconf.FromEnv("MYAPP"),
    ),
)
```