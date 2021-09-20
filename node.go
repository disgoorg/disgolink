package disgolink

type Node interface {
	Lavalink() Lavalink
	Send(d interface{}) error
	Open() error
	Close()
	Name() string
	RestClient() RestClient
	RestURL() string
	Options() *NodeOptions
	Stats() *Stats
}

type NodeOptions struct {
	Name     string
	Host     string
	Port     string
	Password string
	Secure   bool
}
