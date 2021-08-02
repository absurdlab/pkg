package dockertestplus

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"time"
)

// Setup initiates the container creation and boot process according to the spec. If successful, Setup will return
// a Closer function to be invoked after all tests are done, in order to release resources.
func Setup(spec Spec) (closer Closer, err error) {
	var (
		pool     *dockertest.Pool
		resource *dockertest.Resource
	)

	pool, err = dockertest.NewPool(spec.DockerEndpoint)
	if err != nil {
		return
	}

	var hcOpts []func(config *docker.HostConfig)
	if spec.AdditionalConfig != nil {
		hcOpts = append(hcOpts, spec.AdditionalConfig)
	}

	resource, err = pool.RunWithOptions(spec.Options, hcOpts...)
	if err != nil {
		return
	}

	if err = pool.Retry(func() error {
		portMappings := map[string]string{}
		for _, exposedPort := range spec.Options.ExposedPorts {
			portMappings[exposedPort] = resource.GetPort(exposedPort)
		}
		return spec.Ping(portMappings)
	}); err != nil {
		return
	}

	if spec.Lifespan > 0 {
		_ = resource.Expire(uint(spec.Lifespan / time.Second))
	}

	closer = func() error {
		return pool.Purge(resource)
	}

	return
}

// Closer should get invoked after all tests are done.
type Closer func() error
