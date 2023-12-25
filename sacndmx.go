package main

import (
	"fmt"
	"os"
	"sacndmx/hardware"
	"sacndmx/network"

	"github.com/jessevdk/go-flags"
)

var CLIOptions struct {
	IPAddr   string `short:"s" long:"sacn-ip" default:"localhost" description:"Set sACN listener IP address"`
	ListIPs  bool   `short:"n" long:"list-ips" description:"List local IPs"`
	ListDevs bool   `short:"l" long:"list-devices" description:"List output devices for selected output type"`
	Device   string `short:"d" long:"device" default:"" description:"Device serial number to connect to. (default: first encountered device)"`
	Reset    bool   `short:"r" long:"reset-on-timeout" description:"Drop DMX output to zero in case of sACN timeout"`
	Mode     string `short:"t" long:"device-type" default:"opendmx" description:"Output device type. Possible values are opendmx and uart"`
}

var dmx hardware.Hardware
var listener network.Network

func main() {

	// Parsing configuration flags
	_, err := flags.Parse(&CLIOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1) // Exit with code 1 if cli flags are not correct
	}

	// Initializing network listener
	listener, err = network.SACNInit()
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	// Initializing hardware driver
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

	// Processing service flags if any. Listing IPs or devices
	if CLIOptions.ListIPs {
		listener.ListIPs()
		os.Exit(2)
	}

	if CLIOptions.ListDevs {
		if err := dmx.List(); err != nil {
			fmt.Println(err)
		}
		os.Exit(2)
	}

	// Setting up device
	if CLIOptions.Device != "" {
		dmx.SelectDevice(CLIOptions.Device)
	}

	// Saying hello :)
	fmt.Println("sACN-DMX is starting...")

	// Setting up receiver socket
	listener.Bind(CLIOptions.IPAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Connecting to a device
	if err := dmx.Connect(); err == nil {
		fmt.Printf("Using %s (%s)\n", dmx.GetDescription(), dmx.GetSerial())
		defer dmx.Close()
	} else {
		fmt.Printf("Error connecting to device, %s\n", err)
		os.Exit(1)
	}

	// Starting up network listener
	ch := listener.Run()
	fmt.Printf("sACN listener started on %s\n", CLIOptions.IPAddr)

	// Starting up DMX output
	go dmx.Run()
	fmt.Printf("DMX stream started on %s (%s)\n", dmx.GetDescription(), dmx.GetSerial())

	// Main loop
	for {
		// Waiting for a message
		msg := <-ch
		// Setting up channel
		dmx.SetChannel(msg.Channel, msg.Value)
	}

}
