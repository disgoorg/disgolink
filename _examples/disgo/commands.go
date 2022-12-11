package main

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/json"
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
	discord.SlashCommandCreate{
		Name:        "skip",
		Description: "Skips the current song",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "amount",
				Description: "The amount of songs to skip",
				Required:    false,
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "volume",
		Description: "Sets the volume of the player",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "volume",
				Description: "The volume to set",
				Required:    true,
				MaxValue:    json.Ptr(1000),
				MinValue:    json.Ptr(0),
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "seek",
		Description: "Seeks to a specific position in the current song",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "position",
				Description: "The position to seek to",
				Required:    true,
			},
			discord.ApplicationCommandOptionInt{
				Name:        "unit",
				Description: "The unit of the position",
				Required:    false,
				Choices: []discord.ApplicationCommandOptionChoiceInt{
					{
						Name:  "Milliseconds",
						Value: int(lavalink.Millisecond),
					},
					{
						Name:  "Seconds",
						Value: int(lavalink.Second),
					},
					{
						Name:  "Minutes",
						Value: int(lavalink.Minute),
					},
					{
						Name:  "Hours",
						Value: int(lavalink.Hour),
					},
				},
			},
		},
	},
	discord.SlashCommandCreate{
		Name:        "shuffle",
		Description: "Shuffles the current queue",
	},
}

func registerCommands(client bot.Client) {
	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), GuildId, commands); err != nil {
		log.Warn(err)
	}
}
