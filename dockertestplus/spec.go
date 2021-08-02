package dockertestplus

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

// Spec is the configuration for the setup.
type Spec struct {
	DockerEndpoint   string
	Options          *dockertest.RunOptions
	AdditionalConfig func(config *docker.HostConfig)
	Ping             func(portMappings map[string]string) error
	// Lifespan is the absolute lifespan of a container. Docker will try to kill it after this duration.
	Lifespan time.Duration
}
