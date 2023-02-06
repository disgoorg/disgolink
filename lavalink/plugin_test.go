package lavalink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginInfo_MarshalJSON(t *testing.T) {
	pluginInfo := PluginInfo{}
	pluginInfo["test"] = []byte(`{"test": "test"}`)
	pluginInfo["test2"] = []byte(`1.5`)
	data, err := pluginInfo.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `{"test":{"test":"test"},"test2":1.5}`, string(data))
}

func TestPluginInfo_UnmarshalJSON(t *testing.T) {
	pluginInfo := PluginInfo{}
	err := pluginInfo.UnmarshalJSON([]byte(`{"test":{"test":"test"},"test2":1.5}`))
	assert.NoError(t, err)
	assert.Equal(t, PluginInfo{"test": []byte(`{"test":"test"}`), "test2": []byte(`1.5`)}, pluginInfo)
}

func TestPluginInfo_Get(t *testing.T) {
	pluginInfo := PluginInfo{}
	pluginInfo["test"] = []byte(`{"test": "test"}`)
	pluginInfo["test2"] = []byte(`1.5`)

	var (
		test struct {
			Test string `json:"test"`
		}
		test2 float64
		test3 float64
	)

	err := pluginInfo.Get("test", &test)
	assert.NoError(t, err)
	assert.Equal(t, struct {
		Test string `json:"test"`
	}{"test"}, test)

	err = pluginInfo.Get("test2", &test2)
	assert.NoError(t, err)
	assert.Equal(t, 1.5, test2)

	err = pluginInfo.Get("test3", &test3)
	assert.Equal(t, err, ErrPluginDataNotFound("test3"))
	assert.Equal(t, 0.0, test3)
}
