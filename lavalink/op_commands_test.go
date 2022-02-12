package lavalink

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureResumingCommand_MarshalJSON(t *testing.T) {
	data, err := json.Marshal(ConfigureResumingCommand{
		Key:     "test",
		Timeout: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, `{"op":"configureResuming","key":"test","timeout":10}`, string(data))
}
