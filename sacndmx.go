package main

import (
	"fmt"
	"net"
	"os"
	"sacndmx/hardware"

	"github.com/Hundemeier/go-sacn/sacn"
	"github.com/jessevdk/go-flags"
)

var CLIOptions struct {
	IPAddr   string `short:"a" long:"addr" default:"localhost" description:"Listener IP address"`
	ListIPs  bool   `short:"i" long:"list-ips" description:"List local IPs"`
	ListDevs bool   `short:"f" long:"list-devices" description:"List devices"`
	Device   string `short:"d" long:"device" default:"" description:"Device serial number to connect to"`
	Reset    bool   `short:"r" long:"reset-output" description:"Drop DMX output to zero in case of sACN timeout"`
	Mode     string `short:"m" long:"device type" default:"opendmx" description:"Output device type. Possible values are opendmx and uart"`
}

var dmx hardware.Hardware

func main() {

	_, err := flags.Parse(&CLIOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1) // Exit with code 1 if cli flags are not correct
	}

	switch CLIOptions.Mode {
	case "uart":
		dmx, err = hardware.UartInit()
	case "opendmx":
		dmx, err = hardware.EnttecOpenDMXInit()
	default:
		fmt.Printf("unknown device type: %s. Bye.\n", CLIOptions.Mode)
		os.Exit(3)
	}

	if err != nil {
		fmt.Println("Cant access hardware driver")
	}

	if CLIOptions.ListDevs {
		if err := dmx.List(); err != nil {
			fmt.Println(err)
		}
		os.Exit(2)
	}

	if CLIOptions.ListIPs {
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
		}
		os.Exit(2)
	}

	recv, err := sacn.NewReceiverSocket(CLIOptions.IPAddr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("sACN-DMX is starting...")

	if err := dmx.Connect(); err == nil {
		fmt.Printf("Using %s (%s)\n", dmx.GetDescription(), dmx.GetSerial())
		defer dmx.Close()
	} else {
		fmt.Printf("Error connecting to device %s, %s\n", dmx.GetSerial(), err)
		os.Exit(1)
	}

	go dmx.Run()

	recv.SetOnChangeCallback(func(old sacn.DataPacket, newD sacn.DataPacket) {
		data := newD.Data()
		for i := 0; i < len(data); i++ {
			dmx.SetChannel((int16)(i+1), data[i])
		}
	})
	recv.SetTimeoutCallback(func(univ uint16) {
		fmt.Println("Timeout detected on universe", univ)
		// Drop all DMX channels to zero on timeout
		if CLIOptions.Reset {
			for i := 0; i < 512; i++ {
				dmx.SetChannel((int16)(i+1), 0)
			}
		}
	})
	recv.Start()
	fmt.Printf("sACN listener started on %s.\n", CLIOptions.IPAddr)

	select {}
}
