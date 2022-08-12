package lavalink

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/disgoorg/snowflake/v2"
)

var ErrInvalidBotToken = errors.New("invalid bot token")

func UserIDFromBotToken(token string) (snowflake.ID, error) {
	token = strings.TrimPrefix(token, "Bot ")
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return 0, ErrInvalidBotToken
	}
	byteID, err := base64.RawStdEncoding.DecodeString(strs[0])
	if err != nil {
		return 0, err
	}
	return snowflake.Parse(string(byteID))
}
