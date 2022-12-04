package disgolink

import (
	"net/http"
	"time"

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
