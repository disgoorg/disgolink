package main

import "github.com/disgoorg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "shuffle",
		Description: "shuffles the current queue",
	},
	discord.SlashCommandCreate{
		Name:        "filter",
		Description: "applies some filters",
	},
	discord.SlashCommandCreate{
		Name:        "pause",
		Description: "pauses the music",
	},
	discord.SlashCommandCreate{
		Name:        "queue",
		Description: "shows you all tracks in queue",
	},
	discord.SlashCommandCreate{
		Name:        "play",
		Description: "plays some music",
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
