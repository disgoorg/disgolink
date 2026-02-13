package disgolink

import (
	"log/slog"
	"net/http"
	"time"
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

// WithLogger lets you inject your own [slog.Logger].
func WithLogger(logger *slog.Logger) ConfigOpt {
	return func(config *config) {
		config.Logger = logger
	}
}

// WithHTTPClient lets you inject your own http.Client for the client to use for REST requests.
func WithHTTPClient(httpClient *http.Client) ConfigOpt {
	return func(config *config) {
		config.HTTPClient = httpClient
	}
}

// WithListeners adds event listeners to the client.
func WithListeners(listeners ...EventListener) ConfigOpt {
	return func(config *config) {
		config.Listeners = append(config.Listeners, listeners...)
	}
}

// WithListenerFunc adds an event listener function to the client.
func WithListenerFunc[E Event](listenerFunc func(e E)) ConfigOpt {
	return WithListeners(NewListenerFunc(listenerFunc))
}

// WithPlugins adds plugins to the client.
func WithPlugins(plugins ...Plugin) ConfigOpt {
	return func(config *config) {
		config.Plugins = append(config.Plugins, plugins...)
	}
}
