package protoreg

import (
	"sync"

	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"go.uber.org/zap"
)

const (
	pluginName     string = "protoreg"
	rootPluginName string = "grpc"
)

type Configurer interface {
	// UnmarshalKey takes a single key and unmarshal it into a Struct.
	UnmarshalKey(name string, out any) error
	// Has checks if a config section exists.
	Has(name string) bool
	// Experimental returns true if experimental mode is enabled.
	Experimental() bool
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

type Plugin struct {
	mu       *sync.RWMutex
	config   *Config
	log      *zap.Logger
	registry *ProtoRegistry
}

func (p *Plugin) Init(cfg Configurer, log Logger) error {
	const op = errors.Op("protoreg_plugin_init")

	if !cfg.Has(pluginName) {
		return errors.E(errors.Disabled)
	}
	if !cfg.Has(rootPluginName) {
		return errors.E(op, errors.Disabled)
	}

	err := cfg.UnmarshalKey(pluginName, &p.config)
	if err != nil {
		return errors.E(op, err)
	}

	err = p.config.InitDefaults()
	if err != nil {
		return errors.E(op, err)
	}

	p.log = log.NamedLogger(pluginName)
	p.mu = &sync.RWMutex{}

	p.registry, err = p.InitRegistry()

	if err != nil {
		return errors.E(op, err)
	}

	p.log.Info("protoreg initialized")

	return nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*Registry)(nil), p.ProtoRegistry),
	}
}

// ProtoRegistry returns a protobuf registry
func (p *Plugin) ProtoRegistry() *ProtoRegistry {
	return p.registry
}
