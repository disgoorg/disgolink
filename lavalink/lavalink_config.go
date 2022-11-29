package lavalink

import (
	"net/http"
	"time"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
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
	UserID     snowflake.ID
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

func WithUserID(userID snowflake.ID) ConfigOpt {
	return func(config *Config) {
		config.UserID = userID
	}
}

func WithUserIDString(userID string) ConfigOpt {
	parsed, _ := snowflake.Parse(userID)
	return WithUserID(parsed)
}

func WithUserIDFromBotToken(botToken string) ConfigOpt {
	token, _ := UserIDFromBotToken(botToken)
	return WithUserID(token)
}
