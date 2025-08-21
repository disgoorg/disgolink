package disgolink

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/disgoorg/disgolink/v4/lavalink"
)

func defaultConfig() *config {
	return &config{
		Logger: slog.Default(),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type config struct {
	Logger     *slog.Logger
	HTTPClient *http.Client
	Listeners  []EventListener
	Plugins    []Plugin
}

type ConfigOpt func(config *config)

func (c *config) apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
	c.Logger = c.Logger.With(slog.String("name", "disgolink_client"))
}

// WithLogger lets you inject your own Logger implementing log.Logger
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *config) {
		config.Logger = logger
	}
}

func WithHTTPClient(httpClient *http.Client) ConfigOpt {
	return func(config *config) {
		config.HTTPClient = httpClient
	}
}

func WithListeners(listeners ...EventListener) ConfigOpt {
	return func(config *config) {
		config.Listeners = append(config.Listeners, listeners...)
	}
}

func WithListenerFunc[E lavalink.Message](listenerFunc func(p *Player, e E)) ConfigOpt {
	return WithListeners(NewListenerFunc(listenerFunc))
}

func WithPlugins(plugins ...Plugin) ConfigOpt {
	return func(config *config) {
		config.Plugins = append(config.Plugins, plugins...)
	}
}
