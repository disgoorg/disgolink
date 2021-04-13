module github.com/DisgoOrg/disgolink/disgolink/testbot

go 1.16

replace (
	github.com/DisgoOrg/disgolink => ../../
	github.com/DisgoOrg/disgolink/disgolink => ../
)

require (
	github.com/DisgoOrg/disgo v0.1.7-0.20210413004842-8e47aebb5e22
	github.com/DisgoOrg/disgolink v0.0.0-20210412071636-40769c7951dd
	github.com/DisgoOrg/disgolink/disgolink v0.0.0-20210412094129-4268e770cdc4
	github.com/sirupsen/logrus v1.8.1
)
