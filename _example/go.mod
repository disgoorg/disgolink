module github.com/DisgoOrg/disgolink/_example

go 1.17

replace github.com/DisgoOrg/disgolink => ../

require (
	github.com/DisgoOrg/disgo v0.5.12-0.20210920174831-3033f7f639dd
	github.com/DisgoOrg/disgolink v0.2.0
	github.com/DisgoOrg/log v1.1.2
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)
