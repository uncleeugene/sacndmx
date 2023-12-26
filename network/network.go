package network

type dmxChannel struct {
	Channel int16
	Value   byte
}

type Network interface {
	ListIPs() error
	Bind(string) error
	Run() chan dmxChannel
}
