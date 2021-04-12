package api

type Node interface {
	Send(d interface{})
	Open() error
	Name() string
}

type NodeOptions struct {
	Name     string
	Host     string
	Port     int
	Password string
	Secure   bool
}
