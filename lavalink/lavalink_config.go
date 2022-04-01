package lavalink

import (
	"net/http"

	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake"
)

type Config struct {
	Logger     log.Logger
	HTTPClient *http.Client
	UserID     snowflake.Snowflake
	Plugins    []interface{}
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

func WithUserID(userID snowflake.Snowflake) ConfigOpt {
	return func(config *Config) {
		config.UserID = userID
	}
}

func WithUserIDString(userID string) ConfigOpt {
	return WithUserID(snowflake.Snowflake(userID))
}

func WithUserIDFromBotToken(botToken string) ConfigOpt {
	token, _ := UserIDFromBotToken(botToken)
	return WithUserID(token)
}

func WithPlugins(plugins ...interface{}) ConfigOpt {
	return func(config *Config) {
		config.Plugins = append(config.Plugins, plugins...)
	}
}
