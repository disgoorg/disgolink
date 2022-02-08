module github.com/DisgoOrg/disgolink/arikawalink

go 1.17

replace (
	github.com/DisgoOrg/disgolink/lavalink => ../lavalink
)

require (
	github.com/DisgoOrg/disgolink/lavalink v1.1.2
	github.com/DisgoOrg/snowflake v1.0.4
	github.com/diamondburned/arikawa/v3 v3.0.0-rc.4
)

require (
	github.com/DisgoOrg/log v1.1.2 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
)
