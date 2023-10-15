[![Go Reference](https://pkg.go.dev/badge/github.com/disgoorg/disgolink.svg)](https://pkg.go.dev/github.com/disgoorg/disgolink)
[![Go Report](https://goreportcard.com/badge/github.com/disgoorg/disgolink/v3)](https://goreportcard.com/report/github.com/disgoorg/disgolink)
[![Go Version](https://img.shields.io/github/go-mod/go-version/disgoorg/disgolink?filename=go.mod)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/disgolink/blob/master/LICENSE)
[![Disgolink Version](https://img.shields.io/github/v/release/disgoorg/disgolink?label=release)](https://github.com/disgoorg/disgolink/releases/latest)
[![Support Discord](https://discord.com/api/guilds/817327181659111454/widget.png)](https://discord.gg/NFmvZYmZMF)

<img align="right" src="/.github/disgolink.png" width=192 alt="discord gopher">

# DisGoLink

DisGoLink is a [Lavalink](https://github.com/freyacodes/Lavalink) Client written in [Golang](https://golang.org/) which supports the latest Lavalink 4.0.0+ release and the new plugin system. 

While DisGoLink can be used with any [Discord](https://discord.com) Library [DisGo](https://github.com/disgoorg/disgo) is the best fit for it as usage with other Libraries can be a bit annoying due to different [Snowflake](https://github.com/disgoorg/snowflake) implementations.

* [DiscordGo](https://github.com/bwmarrin/discordgo) `string`
* [Arikawa](https://github.com/diamondburned/arikawa) `type Snowflake uint64`
* [Disgord](https://github.com/andersfylling/disgord) `type Snowflake uint64`
* [DisGo](https://github.com/disgoorg/disgo) `type ID uint64`

This Library uses the [Disgo Snowflake](https://github.com/disgoorg/snowflake) package like DisGo

## Getting Started

### Installing

```sh
go get github.com/disgoorg/disgolink/v3/disgolink
```

## Usage

### Setup

First create a new lavalink instance. You can do this either with

```go
import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

var userID = snowflake.ID(1234567890)
lavalinkClient := disgolink.New(userID)
```

You also need to forward the `VOICE_STATE_UPDATE` and `VOICE_SERVER_UPDATE` events to DisGoLink.
Just register an event listener for those events with your library and call `lavalinkClient.OnVoiceStateUpdate` (make sure to only forward your bots voice update event!) and `lavalinkClient.OnVoiceServerUpdate`


For DisGo this would look like this
```go
client, err := disgo.New(Token,
    bot.WithEventListenerFunc(b.onVoiceStateUpdate),
    bot.WithEventListenerFunc(b.onVoiceServerUpdate),
)

func onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
    lavalinkClient.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
}

func onVoiceServerUpdate(event *events.VoiceServerUpdate) {
    lavalinkClient.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}
```

Then you add your lavalink nodes. This directly connects to the nodes and is a blocking call
```go
node, err := lavalinkClient.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:      "test", // a unique node name
		Address:   "localhost:2333",
		Password:  "youshallnotpass",
		Secure:    false, // ws or wss
		SessionID: "", // only needed if you want to resume a previous lavalink session
})
```

after this you can play songs from lavalinks supported sources.

### Loading a track

To play a track you first need to resolve the song. For this you need to call the Lavalink rest `loadtracks` endpoint which returns a result with various track instances. Those tracks can then be played.
```go
query := "ytsearch:Rick Astley - Never Gonna Give You Up"

err := lavalinkClient.BestNode().LoadTracksHandler(context.TODO(), query, lavalink.NewResultHandler(
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

To play a track we first need to connect to the voice channel.
Connecting to a voice channel differs with every lib but here are some quick usages with some
```go
// DisGo
err := client.UpdateVoiceState(context.TODO(), guildID, channelID, false, false)

// DiscordGo
err := session.ChannelVoiceJoinManual(guildID, channelID, false, false)
```

after this you can get/create your player and play the track
```go
player := lavalinkClient.Player("guild_id") // This will either return an existing or new player

var track lavalink.Track // track from result handler before
err := player.Play(track)
```
now audio should start playing

### Listening for events

You can listen for following lavalink events
* `PlayerUpdateMessage` Emitted every x seconds (default 5) with the current player state
* `PlayerPause` Emitted when the player is paused
* `PlayerResume` Emitted when the player is resumed
* `TrackStart` Emitted when a track starts playing
* `TrackEnd` Emitted when a track ends
* `TrackException` Emitted when a track throws an exception
* `TrackStuck` Emitted when a track gets stuck
* `WebsocketClosed` Emitted when the voice gateway connection to lavalink is closed

for this add and event listener for each event to your `Client` instance when you create it or with `Client.AddEventListener`
```go
lavalinkClient := disgolink.New(userID,
    disgolink.WithListenerFunc(onPlayerUpdate),
    disgolink.WithListenerFunc(onPlayerPause),
	disgolink.WithListenerFunc(onPlayerResume),
	disgolink.WithListenerFunc(onTrackStart),
	disgolink.WithListenerFunc(onTrackEnd),
	disgolink.WithListenerFunc(onTrackException),
	disgolink.WithListenerFunc(onTrackStuck),
	disgolink.WithListenerFunc(onWebSocketClosed),
)

func onPlayerUpdate(player disgolink.Player, event lavalink.PlayerUpdateMessage) {
    // do something with the event
}

func onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
    // do something with the event
}

func onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
    // do something with the event
}

func onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
    // do something with the event
}

func onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
    // do something with the event
}

func onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
    // do something with the event
}

func onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
    // do something with the event
}

func onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
    // do something with the event
}
```

### Plugins

Lavalink added [plugins](https://github.com/freyacodes/Lavalink/blob/master/PLUGINS.md) in `v3.5` . DisGoLink exposes a similar API for you to use. With that you can create plugins which require server & client work.
To see what you can do with plugins see [here](https://github.com/disgoorg/disgolink/blob/v2/disgolink/plugin.go)

You register plugins when creating the client instance like this
```go
lavalinkClient := disgolink.New(userID, disgolink.WithPlugins(yourPlugin))
```

Here is a list of plugins(you can pr your own to here):
* [sponsorblock](https://github.com/disgoorg/sponsorblock-plugin) adds payloads and listeners for [Lavalink Sponsorblock-Plugin](https://github.com/Topis-Lavalink-Plugins/Sponsorblock-Plugin)

## Examples

You can find examples under 
* disgo: [_example](https://github.com/disgoorg/disgolink/tree/v2/_examples/disgo)
* discordgo:   [_examples](https://github.com/disgoorg/disgolink/tree/v2/_examples/discordgo)

## Troubleshooting

For help feel free to open an issue or reach out on [Discord](https://discord.gg/NFmvZYmZMF)

## Contributing

Contributions are welcomed but for bigger changes please first reach out via [Discord](https://discord.gg/NFmvZYmZMF) or create an issue to discuss your intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/disgoorg/disgolink/blob/master/LICENSE). See LICENSE for more information.
