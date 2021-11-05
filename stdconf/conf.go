package stdconf

import "github.com/imdario/mergo"

func Parse(options ...Option) (interface{}, error) {
	p := new(parser)
	for _, opt := range options {
		opt(p)
	}

	if p.newFn == nil {
		panic("constructor function is required")
	}

	if err := p.doParse(); err != nil {
		return nil, err
	}

	return p.dest, nil
}

type parser struct {
	dest         interface{}
	newFn        func() interface{}
	mergeConfigs []func(*mergo.Config)
	sourceFns    []func(p *parser) Source
}

func (p *parser) doParse() error {
	if p.dest == nil {
		p.dest = p.newFn()
	}

	for _, each := range p.sourceFns {
		src, err := each(p).Produce()
		if err != nil {
			return err
		}

		if err := mergo.Merge(p.dest, src, p.mergeConfigs...); err != nil {
			return err
		}
	}

	return nil
}

// Option configures the parser
type Option func(p *parser)

// WithNewFunc provides an Option to set the constructor function for the destination config. This function will be
// used by all sources to create the holding structure before merging it onto the base destination.
func WithNewFunc(fn func() interface{}) Option {
	return func(p *parser) {
		if fn != nil {
			p.newFn = fn
		}
	}
}

// WithDestination provides an Option to set the base destination. This option is useful if caller holds a reference
// to the destination in a concrete type, and do not wish to perform an extra type conversion when parser completes.
// Without this option, parser will just use the constructor function to create the base destination.
func WithDestination(dest interface{}) Option {
	return func(p *parser) {
		if dest != nil {
			p.dest = dest
		}
	}
}

// WithMergoOptions provides an Option to customize the merge behaviour. This behaviour affects all merging processes.
func WithMergoOptions(options ...func(*mergo.Config)) Option {
	return func(p *parser) {
		p.mergeConfigs = append(p.mergeConfigs, options...)
	}
}

// WithSources provides an Option to set the configuration sources. The sources will be applied in sequence. By default,
// when the preceding sources have higher priority. This can be changed by provide mergo.WithOverride in WithMergoOptions.
func WithSources(sources ...func(p *parser) Source) Option {
	return func(p *parser) {
		p.sourceFns = append(p.sourceFns, sources...)
	}
}
