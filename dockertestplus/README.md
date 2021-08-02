# dockertestplus

[![.github/workflows/dockertestplus.yaml](https://github.com/absurdlab/pkg/actions/workflows/dockertestplus.yaml/badge.svg)](https://github.com/absurdlab/pkg/actions/workflows/dockertestplus.yaml)

Wrapper around [dockertest](https://github.com/ory/dockertest) to make setup easier.

```shell
go get -u github.com/absurdlab/pkg/dockertestplus
```

## Usage

```go
closer, _ := dockertestplus.Setup(dockertestplus.Spec{
    Options: &dockertest.RunOptions{
        Repository:   "hashicorp/http-echo",
        Tag:          "latest",
        Entrypoint:   []string{"/http-echo", "-text=\"hello world\""},
        ExposedPorts: []string{"5678/tcp"},
    },
    AdditionalConfig: func(config *docker.HostConfig) {
        config.AutoRemove = true
    },
    Ping: func(portMappings map[string]string) error {
        resp, err := http.Get(fmt.Sprintf("http://localhost:%s/", portMappings["5678/tcp"]))
        if err != nil {
            return err
        }

        if resp.StatusCode != http.StatusOK {
            return errors.New("status is not 200")
        }

        return nil
    },
    Lifespan: 10 * time.Minute,
})

defer func() {
	_ = closer()
}
```