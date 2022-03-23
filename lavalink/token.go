package lavalink

import (
	"encoding/base64"
	"strings"

	"github.com/disgoorg/snowflake"
	"github.com/pkg/errors"
)

var ErrInvalidBotToken = errors.New("invalid bot token")

func UserIDFromBotToken(token string) (snowflake.Snowflake, error) {
	token = strings.TrimPrefix(token, "Bot ")
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return "", ErrInvalidBotToken
	}
	byteID, err := base64.StdEncoding.DecodeString(strs[0])
	if err != nil {
		return "", err
	}
	return snowflake.Snowflake(byteID), nil
}
