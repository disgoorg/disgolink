module github.com/DisgoOrg/disgolink/testbot

go 1.16

replace (
	github.com/DisgoOrg/disgolink => ../
	github.com/DisgoOrg/disgolink/disgolink => ../disgolink
)

require (
	github.com/DisgoOrg/disgolink d6c5344
	github.com/DisgoOrg/disgolink/disgolink d6c5344
	github.com/sirupsen/logrus v1.8.1
)
