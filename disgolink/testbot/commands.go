package main

import "github.com/DisgoOrg/disgo/api"

var commands = []api.Command{
	{
		Name:              "play",
		Description:       "plays some music",
		DefaultPermission: true,
		Options: []*api.CommandOption{
			{
				Type:        api.CommandOptionTypeString,
				Name:        "query",
				Description: "what to play",
				Required:    true,
			},
		},
	},
}
