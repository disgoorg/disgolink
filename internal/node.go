package api

type Node struct{
	options NodeOptions

}

type NodeOptions struct{
	host string
	port int
	password string
	secure bool
}
