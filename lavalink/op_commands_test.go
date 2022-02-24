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

func TestFiltersCommand_MarshalJSON(t *testing.T) {
	data, err := json.Marshal(FiltersCommand{
		GuildID: "1234",
		Filters: &DefaultFilters{
			FilterTimescale: &Timescale{
				Speed: 2,
				Pitch: 2,
				Rate:  2,
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, `{"op":"filters","guildId":"1234","timescale":{"speed":2,"pitch":2,"rate":2}}`, string(data))
}
