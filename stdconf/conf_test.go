package stdconf_test

import (
	"github.com/absurdlab/pkg/stdconf"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		env     map[string]string
		options []stdconf.Option
		assert  func(t *testing.T, value interface{}, err error)
	}{
		{
			name: "yaml with env",
			env: map[string]string{
				"TEST_NEST_STRING": "hello world",
			},
			options: []stdconf.Option{
				stdconf.WithNewFunc(newConfig),
				stdconf.WithMergoOptions(mergo.WithOverride),
				stdconf.WithSources(
					stdconf.FromYAMLFile("testdata/config.yaml"),
					stdconf.FromEnv("TEST"),
				),
			},
			assert: func(t *testing.T, value interface{}, err error) {
				if assert.NoError(t, err) && assert.IsType(t, value, &config{}) {
					c := value.(*config)
					assert.Equal(t, "hello", c.String)
					assert.Equal(t, int64(64), c.Int64)
					assert.Equal(t, int(32), c.Int)
					assert.Equal(t, true, c.Bool)
					assert.Equal(t, "hello world", c.Nest.String)
				}
			},
		},
		{
			name: "env with external",
			env: map[string]string{
				"TEST_NEST_STRING": "hello world",
			},
			options: []stdconf.Option{
				stdconf.WithNewFunc(newConfig),
				stdconf.WithMergoOptions(mergo.WithOverride),
				stdconf.WithSources(
					stdconf.FromValue(&config{
						String: "hello",
						Int64:  64,
						Int:    32,
						Bool:   true,
						Nest:   &nested{String: "world"},
					}),
					stdconf.FromEnv("TEST"),
				),
			},
			assert: func(t *testing.T, value interface{}, err error) {
				if assert.NoError(t, err) && assert.IsType(t, value, &config{}) {
					c := value.(*config)
					assert.Equal(t, "hello", c.String)
					assert.Equal(t, int64(64), c.Int64)
					assert.Equal(t, int(32), c.Int)
					assert.Equal(t, true, c.Bool)
					assert.Equal(t, "hello world", c.Nest.String)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if len(c.env) > 0 {
				for k, v := range c.env {
					_ = os.Setenv(k, v)
				}

				defer func() {
					for k := range c.env {
						_ = os.Unsetenv(k)
					}
				}()
			}

			v, err := stdconf.Parse(c.options...)

			c.assert(t, v, err)
		})
	}
}

func newConfig() interface{} {
	return new(config)
}

type config struct {
	String string  `yaml:"string"`
	Int64  int64   `yaml:"int_64"`
	Int    int     `yaml:"int"`
	Bool   bool    `yaml:"bool"`
	Nest   *nested `yaml:"nest,omitempty"`
}

type nested struct {
	String string `yaml:"string"`
}
