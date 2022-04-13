[![Go Reference](https://pkg.go.dev/badge/github.com/disgoorg/disgolink.svg)](https://pkg.go.dev/github.com/disgoorg/disgolink)
[![Go Report](https://goreportcard.com/badge/github.com/disgoorg/disgolink)](https://goreportcard.com/report/github.com/disgoorg/disgolink)
[![Go Version](https://img.shields.io/github/go-mod/go-version/disgoorg/disgolink?filename=lavalink%2Fgo.mod)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/disgolink/blob/master/LICENSE)
[![Disgolink Version](https://img.shields.io/github/v/release/disgoorg/disgolink?label=release)](https://github.com/disgoorg/disgolink/releases/latest)
[![Support Discord](https://discord.com/api/guilds/817327181659111454/widget.png)](https://discord.gg/NFmvZYmZMF)

<img align="right" src="/.github/disgolink.png" width=192 alt="discord gopher">

# DisGolink

DisGolink is a [Lavalink](https://github.com/freyacodes/Lavalink) Client written in [Golang](https://golang.org/) which supports the latest Lavalink 3.4 release and the new plugin system([lavalink dev](https://github.com/freyacodes/Lavalink/tree/dev) only). 

While DisGoLink can be used with any [Discord](https://discord.com) Library [DisGo](https://github.com/disgoorg/disgo) is the best fit for it as usage with other Libraries can be a bit annoying due to different [Snowflake](https://github.com/disgoorg/snowflake) implementations.

* [DiscordGo](https://github.com/bwmarrin/discordgo) `string`
* [Arikawa](https://github.com/diamondburned/arikawa) `type Snowflake uint64`
* [Disgord](https://github.com/andersfylling/disgord) `type Snowflake uint64`
* [DisGo](https://github.com/disgoorg/disgo) `type Snowflake string`

This Library uses the [Disgo Snowflake](https://github.com/disgoorg/snowflake) package like DisGo

## Getting Started

### Installing

For `DisGo` and `DiscordGo` there is a sub module which simplifies the usage of DisGolink a bit. You can skip those if you want and directly get the Lavalink Client via

```sh
go get github.com/disgoorg/disgolink/lavalink
```

or you can get the library specific packages via

#### DisGo
```sh
go get github.com/disgoorg/disgolink/disgolink
```

#### DiscordGo
```sh
go get github.com/disgoorg/disgolink/dgolink
```

## Usage

### Setup

First create a new lavalink instance. You can do this either with

```go
import "github.com/disgoorg/disgolink/lavalink"

link := lavalink.New(lavalink.WithUserID("user_id_here"))
```

or with the library specific packages

#### DisGo
```go
import "github.com/disgoorg/disgolink/disgolink"

link := disgolink.New(disgo)
```

#### DiscordGo
```go
import "github.com/disgoorg/disgolink/dgolink"

link := dgolink.New(session)
```

then you add your lavalink nodes. This directly connects to the nodes and is a blocking call
```go
node, err := link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:        "test", // a unique node name
		Host:        "localhost",
		Port:        "2333",
		Password:    "youshallnotpass",
		Secure:      false, // ws or wss
		ResumingKey: "", // only needed if you want to resume a lavalink session
})
```

after this you can play songs from lavalinks supported sources.

### Loading a track

To play a track you first need to resolve the song. For this you need to call the Lavalink rest `loadtracks` endpoint which returns a result with various track instances. Those tracks can then be played.
```go
query := "ytsearch:Rick Astley - Never Gonna Give You Up"

err := link.BestRestClient().LoadItemHandler(context.TODO(), query, lavalink.NewResultHandler(
		func(track lavalink.AudioTrack) {
				// Loaded a single track
		},
		func(playlist lavalink.AudioPlaylist) {
				// Loaded a playlist
		},
		func(tracks []lavalink.AudioTrack) {
				// Loaded a search result
		},
		func() {
				// nothing matching the query found
		},
		func(ex lavalink.FriendlyException) {
				// something went wrong while loading the track
		},
))
```

### Playing a track

To play a track we first need to to connect to the voice channel.
Connecting to a voice channel differs with every lib but here are some quick usages with some
```go
// DisGo
err := client.Connect(context.TODO(), guildID, channelID)

// DiscordGo
err := session.ChannelVoiceJoinManual(guildID, channelID, false, false)
```

after this you can get/create your player and play the track
```go
player := link.Player("guild_id") // This will either return an existing or new player

var track lavalink.channelID // track from result handler before
err := player.Play(track)
```
now audio should start playing

### Listening for events

You can listen for following lavalink events
* `PlayerUpdate`
* `PlayerPause`
* `PlayerResume`
* `TrackStart`
* `TrackEnd`
* `TrackException`
* `TrackStuck`
* `WebsocketClosed`

for this implement the [`PlayerEventListener`](https://github.com/disgoorg/disgolink/blob/master/lavalink/player_listener.go) interface. 
To listen to only a few events you can optionally embed the `PlayerEventAdapter` struct which has empty dummy methods.

After implementing the interface you can add the listener to the player and your methods should start getting called
```go

type EventListener struct {
    lavalink.PlayerEventAdapter
}
func (l *EventListener) OnTrackStart(player Player, track AudioTrack)                                  {}
func (l *EventListener) OnTrackEnd(player Player, track AudioTrack, endReason AudioTrackEndReason)     {}
func (l *EventListener) OnTrackException(player Player, track AudioTrack, exception FriendlyException) {}

player.AddListener(&EventListener{})
```

### Plugins

Lavalink added plugins on the [dev branch](https://github.com/freyacodes/Lavalink/blob/dev/PLUGINS.md). DisGolink exposes a similar API for you to use. With that you can create plugins which require server & client work.
To see what you can do with plugins see [here](https://github.com/disgoorg/disgolink/blob/master/lavalink/plugin.go)

You register plugins when creating the link instance like this
```go
link := lavalink.New(lavalink.WithUserID("user_id_here"), lavalink.WithPlugins(yourPlugin))

// DisGo
link := disgolink.New(client, lavalink.WithPlugins(yourPlugin)) 

// DiscordGo
link := dgolink.New(session, lavalink.WithPlugins(yourPlugin))
```

Here is a list of plugins(you can pr your own to here):
* [source-extensions](https://github.com/disgoorg/source-extensions-plugin) adds source track encoder & decoder for [Lavalink Topis-Source-Managers-Plugin](https://github.com/Topis-Lavalink-Plugins/Topis-Source-Managers-Plugin)
* [sponsorblock](https://github.com/disgoorg/sponsorblock-plugin) adds payloads and listeners for [Lavalink Sponsorblock-Plugin](https://github.com/Topis-Lavalink-Plugins/Sponsorblock-Plugin)

## Examples

You can find examples under 
* lavalink:  [_example](https://github.com/disgoorg/disgolink/tree/master/_example)
* disgolink: [_example](https://github.com/disgoorg/disgolink/tree/master/disgolink/_example)
* dgolink:   [_examples](https://github.com/disgoorg/disgolink/tree/master/dgolink/_example)

## Troubleshooting

For help feel free to open an issues or reach out on [Discord](https://discord.gg/NFmvZYmZMF)

## Contributing

Contributions are welcomed but for bigger changes please first reach out via [Discord](https://discord.gg/NFmvZYmZMF) or create an issue to discuss your intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/disgolink/blob/master/LICENSE). See LICENSE for more information.

