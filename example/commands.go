package main

import dapi "github.com/DisgoOrg/disgo/api"

var commands = []dapi.CommandCreate{
	{
		Name:        "shuffle",
		Description: "shuffles the current queue",
	},
	{
		Name:        "filter",
		Description: "applies some filters",
	},
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
		Options: []dapi.CommandOption{
			{
				Type:        dapi.CommandOptionTypeString,
				Name:        "query",
				Description: "what to play",
				Required:    true,
			},
			{
				Type:        dapi.CommandOptionTypeString,
				Name:        "search-provider",
				Description: "where to search",
				Choices: []dapi.OptionChoice{
					{Name: "YouTube", Value: "yt"},
					{Name: "YouTube Music", Value: "ytm"},
					{Name: "SoundCloud", Value: "sc"},
				},
			},
		},
	},
}
