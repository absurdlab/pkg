package dockertestplus

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

// Spec is the configuration for the setup.
type Spec struct {
	// DockerEndpoint is the endpoint to the Docker. In most cases, Docker is installed locally,
	// and this field can be left empty.
	DockerEndpoint string
	// Options is the container configuration.
	Options *dockertest.RunOptions
	// AdditionalConfig allows deep customization of the final docker config. For example, to enable AutoRemove feature.
	AdditionalConfig func(config *docker.HostConfig)
	// Ping is a health check function asserting the created resource is live and ready. It is supplied with a map
	// of port mappings, indexed with Options.ExposedPorts.
	Ping func(portMappings map[string]string) error
	// Lifespan is the absolute lifespan of a container. Docker will try to kill it after this duration.
	Lifespan time.Duration
}
