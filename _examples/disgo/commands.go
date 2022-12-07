package main

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/log"
)

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "play",
		Description: "Plays a song",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "identifier",
				Description: "The song link or search query",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "source",
				Description: "The source to search on",
				Required:    false,
				Choices: []discord.ApplicationCommandOptionChoiceString{
					{
						Name:  "YouTube",
						Value: string(lavalink.SearchTypeYoutube),
					},
					{
						Name:  "YouTube Music",
						Value: string(lavalink.SearchTypeYoutubeMusic),
					},
					{
						Name:  "SoundCloud",
						Value: string(lavalink.SearchTypeSoundCloud),
					},
					{
						Name:  "Deezer",
						Value: "dzsearch",
					},
					{
						Name:  "Deezer ISRC",
						Value: "dzisrc",
					},
					{
						Name:  "Spotify",
						Value: "spsearch",
					},
					{
						Name:  "AppleMusic",
						Value: "amsearch",
					},
				},
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "pause",
		Description: "Pauses the current song",
	},
	discord.SlashCommandCreate{
		Name:        "now-playing",
		Description: "Shows the current playing song",
	},
	discord.SlashCommandCreate{
		Name:        "stop",
		Description: "Stops the current song and stops the player",
	},
	discord.SlashCommandCreate{
		Name:        "players",
		Description: "Shows all active players",
	},
}

func registerCommands(client bot.Client) {
	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), GUILD_ID, commands); err != nil {
		log.Warn(err)
	}
}
