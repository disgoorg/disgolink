package main

import "github.com/disgoorg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		CommandName: "shuffle",
		Description: "shuffles the current queue",
	},
	discord.SlashCommandCreate{
		CommandName: "filter",
		Description: "applies some filters",
	},
	discord.SlashCommandCreate{
		CommandName: "pause",
		Description: "pauses the music",
	},
	discord.SlashCommandCreate{
		CommandName: "queue",
		Description: "shows you all tracks in queue",
	},
	discord.SlashCommandCreate{
		CommandName: "play",
		Description: "plays some music",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				OptionName:  "query",
				Description: "what to play",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				OptionName:  "search-provider",
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
