package main

import "github.com/DisgoOrg/disgo/api"

var commands = []api.CommandCreate{
	{
		Name:        "pause",
		Description: "pauses the music",
	},
	{
		Name:        "queue",
		Description: "shows you all tracks in queue",
	},
	{
		Name:        "play",
		Description: "plays some music",
		Options: []api.CommandOption{
			{
				Type:        api.CommandOptionTypeString,
				Name:        "query",
				Description: "what to play",
				Required:    true,
			},
			{
				Type:        api.CommandOptionTypeString,
				Name:        "search-provider",
				Description: "where to search",
				Choices: []api.OptionChoice{
					{Name: "YouTube", Value: "yt"},
					{Name: "YouTube Music", Value: "ytm"},
					{Name: "SoundCloud", Value: "sc"},
				},
			},
		},
	},
}