module github.com/DisgoOrg/disgolink/disgolink

go 1.17

replace (
	github.com/DisgoOrg/disgo => ../../disgo
	github.com/DisgoOrg/disgolink => ../
)

require (
	github.com/DisgoOrg/disgo v0.6.12
	github.com/DisgoOrg/disgolink v1.0.1-0.20220113110532-5b6f72beb7fe
)

require (
	github.com/DisgoOrg/log v1.1.2 // indirect
	github.com/DisgoOrg/snowflake v1.0.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b // indirect
)
