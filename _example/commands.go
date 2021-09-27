package main

import "github.com/DisgoOrg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	{
		Name:              "shuffle",
		Description:       "shuffles the current queue",
		DefaultPermission: true,
	},
	{
		Name:              "filter",
		Description:       "applies some filters",
		DefaultPermission: true,
	},
	{
		Name:              "pause",
		Description:       "pauses the music",
		DefaultPermission: true,
	},
	{
		Name:              "queue",
		Description:       "shows you all tracks in queue",
		DefaultPermission: true,
	},
	{
		Name:              "play",
		Description:       "plays some music",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.ApplicationCommandOptionTypeString,
				Name:        "query",
				Description: "what to play",
				Required:    true,
			},
			{
				Type:        discord.ApplicationCommandOptionTypeString,
				Name:        "search-provider",
				Description: "where to search",
				Choices: []discord.ApplicationCommandOptionChoice{
					{Name: "YouTube", Value: "yt"},
					{Name: "YouTube Music", Value: "ytm"},
					{Name: "SoundCloud", Value: "sc"},
				},
			},
		},
	},
}
