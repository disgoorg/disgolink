package main

import "github.com/DisgoOrg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
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
		Options: []discord.ApplicationCommandOption{
			{
				Type:        discord.CommandOptionTypeString,
				Name:        "query",
				Description: "what to play",
				Required:    true,
			},
			{
				Type:        discord.CommandOptionTypeString,
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
