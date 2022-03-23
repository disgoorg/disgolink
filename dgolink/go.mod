module github.com/disgoorg/disgolink/dgolink

go 1.18

replace github.com/disgoorg/disgolink/lavalink => ../lavalink

require (
	github.com/bwmarrin/discordgo v0.23.3-0.20211228023845-29269347e820
	github.com/disgoorg/disgolink/lavalink v1.4.1
	github.com/disgoorg/snowflake v1.1.0
)

require (
	github.com/disgoorg/log v1.2.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)
