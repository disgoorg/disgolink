module github.com/DisgoOrg/disgolink/_example

go 1.17

replace github.com/DisgoOrg/disgolink/lavalink => ../lavalink

require (
	github.com/DisgoOrg/disgolink/lavalink v0.2.5-0.20220112205449-450f387bf713
	github.com/DisgoOrg/snowflake v1.0.3
)

require (
	github.com/DisgoOrg/log v1.1.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
