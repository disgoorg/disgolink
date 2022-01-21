package lavalink

import (
	"net/http"

	"github.com/DisgoOrg/log"
)

type Config struct {
	Logger     log.Logger
	HTTPClient *http.Client
	UserID     string
	Plugins    []interface{}
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger lets you inject your own logger implementing log.Logger
//goland:noinspection GoUnusedExportedFunction
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithHTTPClient(httpClient *http.Client) ConfigOpt {
	return func(config *Config) {
		config.HTTPClient = httpClient
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithUserID(userID string) ConfigOpt {
	return func(config *Config) {
		config.UserID = userID
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithPlugins(plugins ...interface{}) ConfigOpt {
	return func(config *Config) {
		config.Plugins = append(config.Plugins, plugins...)
	}
}
