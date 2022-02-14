package lavalink

import (
	"net/http"

	"github.com/DisgoOrg/log"
	"github.com/DisgoOrg/snowflake"
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
func WithUserID(userID snowflake.Snowflake) ConfigOpt {
	return func(config *Config) {
		config.UserID = userID
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithUserIDString(userID string) ConfigOpt {
	return WithUserID(snowflake.Snowflake(userID))
}

//goland:noinspection GoUnusedExportedFunction
func WithUserIDFromBotToken(botToken string) ConfigOpt {
	token, _ := UserIDFromBotToken(botToken)
	return WithUserID(token)
}

//goland:noinspection GoUnusedExportedFunction
func WithPlugins(plugins ...interface{}) ConfigOpt {
	return func(config *Config) {
		config.Plugins = append(config.Plugins, plugins...)
	}
}
