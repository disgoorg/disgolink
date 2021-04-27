module github.com/DisgoOrg/disgolink/disgo/example

go 1.16

replace (
	github.com/DisgoOrg/disgolink => ../../
	github.com/DisgoOrg/disgolink/disgo => ../
)

require (
	github.com/DisgoOrg/disgo v0.3.2
	github.com/DisgoOrg/disgolink v0.0.0-20210412071636-40769c7951dd
	github.com/DisgoOrg/disgolink/disgo v0.0.0-20210412094129-4268e770cdc4
	github.com/sirupsen/logrus v1.8.1
)
