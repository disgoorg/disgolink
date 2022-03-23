package main

import "github.com/disgoorg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		CommandName:       "shuffle",
		Description:       "shuffles the current queue",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		CommandName:       "filter",
		Description:       "applies some filters",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		CommandName:       "pause",
		Description:       "pauses the music",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		CommandName:       "queue",
		Description:       "shows you all tracks in queue",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		CommandName:       "play",
		Description:       "plays some music",
		DefaultPermission: true,
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "query",
				Description: "what to play",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "search-provider",
				Description: "where to search",
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{Name: "YouTube", Value: "yt"},
					{Name: "YouTube Music", Value: "ytm"},
					{Name: "SoundCloud", Value: "sc"},
				},
			},
		},
	},
}
