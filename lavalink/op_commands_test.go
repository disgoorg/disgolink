package lavalink

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigureResumingCommand_MarshalJSON(t *testing.T) {
	data, err := json.Marshal(ConfigureResumingCommand{
		Key:     "test",
		Timeout: 10 * time.Second,
	})
	assert.NoError(t, err)
	assert.Equal(t, `{"op":"configureResuming","timeout":10000,"key":"test"}`, string(data))
}
