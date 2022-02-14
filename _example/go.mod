module github.com/DisgoOrg/disgolink/_example

go 1.17

replace github.com/DisgoOrg/disgolink/lavalink => ../lavalink

require (
	github.com/DisgoOrg/disgolink/lavalink v1.3.3
	github.com/DisgoOrg/snowflake v1.0.4
)

require (
	github.com/DisgoOrg/log v1.1.3 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
