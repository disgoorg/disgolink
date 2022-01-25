[![Go Reference](https://pkg.go.dev/badge/github.com/DisgoOrg/disgolink.svg)](https://pkg.go.dev/github.com/DisgoOrg/disgolink)
[![Go Report](https://goreportcard.com/badge/github.com/DisgoOrg/disgolink)](https://goreportcard.com/report/github.com/DisgoOrg/disgolink)
[![Go Version](https://img.shields.io/github/go-mod/go-version/DisgoOrg/disgolink?filename=lavalink%2Fgo.mod)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/disgolink/blob/master/LICENSE)
[![Disgo Version](https://img.shields.io/github/v/release/DisgoOrg/disgolink?label=release)](https://github.com/DisgoOrg/disgolink/releases/latest)
[![Disgo Discord](https://discord.com/api/guilds/817327181659111454/widget.png)](https://discord.gg/NFmvZYmZMF)

<img align="right" src="/.github/disgolink.png" width=192 alt="discord gopher">

# disgolink

disgolink is a [Lavalink](https://github.com/freyacodes/Lavalink) Client which supports the latest Lavalink 3.4 release

## Getting Started

### Installing

There are 3 packages depending on which go lib you use get a different package

#### lavalink(non specific implementation)

```sh
go get github.com/DisgoOrg/disgolink/lavalink
```

#### disgolink([disgo](https://github.com/DisgoOrg/disgo) implementation)

```sh
go get github.com/DisgoOrg/disgolink/disgolink
```

#### dgolink([discordgo](https://github.com/bwmarrin/discordgo) implementation)

```sh
go get github.com/DisgoOrg/disgolink/dgolink
```

### Building a Lavalink instance

#### lavalink

```go
import "github.com/DisgoOrg/disgolink/lavalink"

link := lavalink.New(lavalink.WithUserID("user_id_here"))
```

#### disgolink

```go
import "github.com/DisgoOrg/disgolink/dgolink"

link := dgolink.New(session)
```

#### dgolink

```go
import "github.com/DisgoOrg/disgolink/lavalink"

link := disgolink.New(disgo)
```

## Examples

You can find examples under 
* lavalink:  [_example](https://github.com/DisgoOrg/disgolink/tree/master/_example)
* disgolink: [_example](https://github.com/DisgoOrg/disgolink/tree/master/disgolink/_example)
* dgolink:   [_examples](https://github.com/DisgoOrg/disgolink/tree/master/dgolink/_example)

## Troubleshooting

For help feel free to open an issues or reach out on [Discord](https://discord.gg/NFmvZYmZMF)

## Contributing

Contributions are welcomed but for bigger changes please first reach out via [Discord](https://discord.gg/NFmvZYmZMF) or create an issue to discuss your intentions and ideas.

## License

Distributed under the [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/DisgoOrg/disgolink/blob/master/LICENSE). See LICENSE for more information.

