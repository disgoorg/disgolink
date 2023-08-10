package disgolink

import (
	"net/http"
	"time"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/log"
)

func DefaultConfig() *Config {
	return &Config{
		Logger:     log.Default(),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type Config struct {
	Logger     log.Logger
	HTTPClient *http.Client
	Listeners  []EventListener
	Plugins    []Plugin
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger lets you inject your own logger implementing log.Logger
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

func WithHTTPClient(httpClient *http.Client) ConfigOpt {
	return func(config *Config) {
		config.HTTPClient = httpClient
	}
}

func WithListeners(listeners ...EventListener) ConfigOpt {
	return func(config *Config) {
		config.Listeners = append(config.Listeners, listeners...)
	}
}

func WithListenerFunc[E lavalink.Message](listenerFunc func(p Player, e E)) ConfigOpt {
	return WithListeners(NewListenerFunc(listenerFunc))
}

func WithPlugins(plugins ...Plugin) ConfigOpt {
	return func(config *Config) {
		config.Plugins = append(config.Plugins, plugins...)
	}
}
