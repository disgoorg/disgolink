package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/log"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "play",
		Description: "Plays a song",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "identifier",
				Description: "The song link or search query",
				Required:    true,
			},
		},
	},
	{
		Name:        "pause",
		Description: "Pauses the current song",
	},
	{
		Name:        "now-playing",
		Description: "Shows the current playing song",
	},
	{
		Name:        "stop",
		Description: "Stops the current song and stops the player",
	},
	{
		Name:        "players",
		Description: "Shows all active players",
	},
	{
		Name:        "shuffle",
		Description: "Shuffles the current queue",
	},
	{
		Name:        "queue",
		Description: "Shows the current queue",
	},
	{
		Name:        "clear-queue",
		Description: "Clears the current queue",
	},
	{
		Name:        "queue-type",
		Description: "Sets the queue type",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "type",
				Description: "The queue type",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "default",
						Value: "default",
					},
					{
						Name:  "repeat-track",
						Value: "repeat-track",
					},
					{
						Name:  "repeat-queue",
						Value: "repeat-queue",
					},
				},
			},
		},
	},
}

func registerCommands(s *discordgo.Session) {
	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildId, commands); err != nil {
		log.Warn(err)
	}
}
