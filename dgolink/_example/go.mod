module github.com/DisgoOrg/disgolink/dgolink/_example

go 1.17

replace (
	github.com/DisgoOrg/disgolink => ../../
	github.com/DisgoOrg/disgolink/dgolink => ../
)

require (
	github.com/DisgoOrg/disgolink v0.2.5-0.20220112205449-450f387bf713
	github.com/DisgoOrg/disgolink/dgolink v0.0.0-20220112214621-4d364f10b07e
	github.com/bwmarrin/discordgo v0.23.3-0.20211228023845-29269347e820
)

require (
	github.com/DisgoOrg/log v1.1.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)
