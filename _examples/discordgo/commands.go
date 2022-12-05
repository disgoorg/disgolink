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
}

func registerCommands(s *discordgo.Session) {
	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, GUILD_ID, commands); err != nil {
		log.Warn(err)
	}
}
