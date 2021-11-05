package stdconf

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

// Source abstracts where a configuration comes from.
type Source interface {
	// Produce returns a structure parsed from the source. The structure
	// returned must be of the same type as the base structure.
	Produce() (interface{}, error)
}

// FromValue returns a Source which produces the supplied value as is. It is useful when supplying default values, or
// a structure injected from external sources (i.e. command line args)
func FromValue(value interface{}) func(p *parser) Source {
	return func(p *parser) Source {
		return &valueSource{value: value}
	}
}

// FromJSONFile returns a Source that reads configuration from a JSON file.
func FromJSONFile(file string) func(p *parser) Source {
	return func(p *parser) Source {
		return &jsonSource{
			readFn: func() (io.Reader, error) {
				return os.Open(file)
			},
			newFn: p.newFn,
		}
	}
}

// FromJSONString returns a Source that reads configuration from a JSON string.
func FromJSONString(value string) func(p *parser) Source {
	return func(p *parser) Source {
		return &jsonSource{
			readFn: func() (io.Reader, error) {
				return strings.NewReader(value), nil
			},
			newFn: p.newFn,
		}
	}
}

// FromYAMLFile returns a Source that reads configuration from a YAML string.
func FromYAMLFile(file string) func(p *parser) Source {
	return func(p *parser) Source {
		return &yamlSource{
			readFn: func() (io.Reader, error) {
				return os.Open(file)
			},
			newFn: p.newFn,
		}
	}
}

// FromYAMLString returns a Source that reads configuration from a YAML string.
func FromYAMLString(value string) func(p *parser) Source {
	return func(p *parser) Source {
		return &yamlSource{
			readFn: func() (io.Reader, error) {
				return strings.NewReader(value), nil
			},
			newFn: p.newFn,
		}
	}
}

// FromEnv returns a Source that reads configuration from environment variable.
func FromEnv(prefix string) func(p *parser) Source {
	return func(p *parser) Source {
		return &envSource{
			prefix: prefix,
			newFn:  p.newFn,
		}
	}
}

type valueSource struct {
	value interface{}
}

func (v *valueSource) Produce() (interface{}, error) {
	return v.value, nil
}

type jsonSource struct {
	readFn func() (io.Reader, error)
	newFn  func() interface{}
}

func (j *jsonSource) Produce() (interface{}, error) {
	reader, err := j.readFn()
	if err != nil {
		return nil, err
	}

	dest := j.newFn()
	if err := json.NewDecoder(reader).Decode(dest); err != nil {
		return nil, err
	}

	return dest, nil
}

type yamlSource struct {
	readFn func() (io.Reader, error)
	newFn  func() interface{}
}

func (y *yamlSource) Produce() (interface{}, error) {
	reader, err := y.readFn()
	if err != nil {
		return nil, err
	}

	dest := y.newFn()
	if err := yaml.NewDecoder(reader).Decode(dest); err != nil {
		return nil, err
	}

	return dest, nil
}

type envSource struct {
	prefix string
	newFn  func() interface{}
}

func (e *envSource) Produce() (interface{}, error) {
	dest := e.newFn()

	if err := envconfig.Process(e.prefix, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
