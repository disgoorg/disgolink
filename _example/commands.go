package main

import "github.com/DisgoOrg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:              "shuffle",
		Description:       "shuffles the current queue",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		Name:              "filter",
		Description:       "applies some filters",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		Name:              "pause",
		Description:       "pauses the music",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		Name:              "queue",
		Description:       "shows you all tracks in queue",
		DefaultPermission: true,
	},
	discord.SlashCommandCreate{
		Name:              "play",
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
			discord.ApplicationCommandOptionBool{
				Name:        "skip-segments",
				Description: "skip sponsorblock segments",
			},
		},
	},
}
