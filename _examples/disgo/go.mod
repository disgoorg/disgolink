module github.com/disgoorg/disgolink/v2/_examples/disgo

go 1.18

replace github.com/disgoorg/disgolink/v2 => ../../

require (
	github.com/disgoorg/disgo v0.14.0
	github.com/disgoorg/disgolink/v2 v2.0.0
	github.com/disgoorg/json v1.0.0
	github.com/disgoorg/log v1.2.0
	github.com/disgoorg/snowflake/v2 v2.0.1
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b // indirect
	golang.org/x/exp v0.0.0-20221126150942-6ab00d035af9 // indirect
)
