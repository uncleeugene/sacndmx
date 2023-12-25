package network

import (
	"fmt"
	"net"

	"github.com/Hundemeier/go-sacn/sacn"
)

type sACN struct {
	Socket *sacn.ReceiverSocket
	Out    chan dmxChannel
}

func SACNInit() (*sACN, error) {
	n := new(sACN)
	n.Out = make(chan dmxChannel)
	return n, nil
}

func (n *sACN) ListIPs() error {
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

func (n *sACN) Bind(addr string) error {
	recv, err := sacn.NewReceiverSocket(addr, nil)
	if err == nil {
		n.Socket = recv
		n.Socket.SetOnChangeCallback(n.OnChangeCallback)
		n.Socket.SetTimeoutCallback(n.OnTimeoutCallback)
		return nil
	} else {
		return err
	}
}

func (n *sACN) Output() chan dmxChannel {
	return n.Out
}

func (n *sACN) Run() chan dmxChannel {
	n.Socket.Start()
	return n.Out
}

func (n *sACN) OnChangeCallback(old sacn.DataPacket, newD sacn.DataPacket) {
	data := newD.Data()
	for i := 0; i < len(data); i++ {
		ch := dmxChannel{
			(int16)(i + 1),
			data[i],
		}
		n.Out <- ch
	}
}

func (n *sACN) OnTimeoutCallback(univ uint16) {
	fmt.Println("Timeout detected on universe", univ)
	// Drop all DMX channels to zero on timeout
	if false {
		for i := 0; i < 512; i++ {
			n.Out <- dmxChannel{
				(int16)(i),
				0,
			}
		}
	}
}
