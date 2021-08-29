module github.com/DisgoOrg/disgolink/_example

go 1.17

replace (
	github.com/DisgoOrg/disgolink => ../
	github.com/DisgoOrg/disgo => ../../disgo
)

require (
	github.com/DisgoOrg/disgo v0.5.7
	github.com/DisgoOrg/disgolink v0.2.0
	github.com/DisgoOrg/log v1.1.0
)
