package network

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet"
	"github.com/jsimonetti/go-artnet/packet"
	"github.com/jsimonetti/go-artnet/packet/code"
)

type artNet struct {
	node *artnet.Node
	Out  chan dmxChannel
}

func ArtNetInit() (*artNet, error) {
	a := new(artNet)
	a.Out = make(chan dmxChannel)
	return a, nil
}

func (a *artNet) ListIPs() error {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					fmt.Printf("Local addr: %s\n", ipnet.IP.String())
				}
			}
		}
		return nil
	} else {
		return err
	}
}

func (n *artNet) Bind(addr string) error {
	if ip, _, err := net.ParseCIDR(addr); err == nil {
		log := artnet.NewDefaultLogger()
		n.node = artnet.NewNode("sacndmx", code.StNode, ip, log)
		n.node.RegisterCallback(code.OpDMX, n.OnChangeCallback)
		return nil
	} else {
		return err
	}
}

func (n *artNet) Run() chan dmxChannel {
	n.node.Start()
	return n.Out
}

func (n *artNet) OnChangeCallback(p packet.ArtNetPacket) {
	fmt.Println("Got packet")
	dmx, ok := p.(*packet.ArtDMXPacket)
	if !ok {
		fmt.Println("Invalid DMX Packet")
	} else {
		stream := dmx.Data
		for i := range stream {
			ch := dmxChannel{
				(int16)(i + 1),
				stream[i],
			}
			n.Out <- ch
		}
	}

}
