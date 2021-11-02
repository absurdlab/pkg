package dockertestplus_test

import (
	"errors"
	"fmt"
	"github.com/absurdlab/pkg/dockertestplus"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"net/http"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	closer, err := dockertestplus.Setup(dockertestplus.Spec{
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
	if err != nil {
		t.Error(err)
	}

	defer func() {
		_ = closer()
	}()
}
