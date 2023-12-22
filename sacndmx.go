package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Hundemeier/go-sacn/sacn"
	"github.com/jessevdk/go-flags"
	ftdi "github.com/uncleeugene/goftdi"
)

type enttecOpenDMX struct {
	Device  *ftdi.Device
	channel []byte
	buffer  []byte
}

var lineProperties = ftdi.LineProperties{
	Bits:     ftdi.BITS_8,
	StopBits: ftdi.STOP_2,
	Parity:   ftdi.NONE,
}

func enttecOpenDMXConnect(d ftdi.DeviceInfo) (enttecOpenDMX, error) {
	var device enttecOpenDMX
	dev, err := ftdi.Open(d)
	if err == nil {
		device.Device = dev
		device.channel = make([]byte, 512)
		device.buffer = make([]byte, 513)
	}
	device.Device.SetLineProperty(lineProperties)
	device.Device.SetBaudRate(250000)
	device.Device.Purge()
	return device, err
}

func (d *enttecOpenDMX) Close() error {
	return d.Device.Close()
}

func (d *enttecOpenDMX) SetChannel(index int16, data byte) error {
	d.channel[index-1] = data

	return nil
}

func (d *enttecOpenDMX) Render() error {
	for i := 0; i < 512; i++ {
		d.buffer[i+1] = d.channel[i]
	}
	d.Device.SetBreakOn(lineProperties)
	d.Device.SetBreakOff(lineProperties)
	if _, err := d.Device.Write(d.buffer); err != nil {
		return err
	}
	return nil
}

func (d *enttecOpenDMX) Run() {
	for {
		if err := d.Render(); err != nil {
			log.Fatalf("Error rendering DMX: %s\n", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

var CLIOptions struct {
	DumpTOML bool   `short:"s" long:"showconfig" description:"Dump configuration and exit"`
	Config   string `short:"c" long:"config" default:"sacndmx.toml" description:"Configuration file path"`
	IPAddr   string `short:"a" long:"addr" default:"localhost" description:"Listener IP address"`
	ListIPs  bool   `short:"i" long:"list-ips" description:"List local IPs"`
	ListDevs bool   `short:"f" long:"list-devices" description:"List devices"`
	Device   string `short:"d" long:"device" default:"" description:"Device serial number to connect to"`
	Reset    bool   `short:"r" long:"reset-output" description:"Drop DMX output to zero in case of sACN timeout"`
	DevType  string `short:"t" long:"device-type" default:"opendmx" description:"Device type. Not implemented yet."`
}

func main() {
	var exitFlag bool

	_, err := flags.Parse(&CLIOptions)
	if err != nil {
		log.Println(err)
		os.Exit(1) // Exit with code 1 if cli flags are not correct
	}

	dl, err := ftdi.GetDeviceList()
	if err == nil {
		if CLIOptions.ListDevs {
			for i := 0; i < len(dl); i++ {
				fmt.Printf("Device %d: S/N %s, Desc: \"%s\"\n", i, dl[i].SerialNumber, dl[i].Description)
			}
			exitFlag = true
		}
	} else {
		log.Fatal(err)
		exitFlag = true
	}

	recv, err := sacn.NewReceiverSocket(CLIOptions.IPAddr, nil)
	if err != nil {
		log.Fatal(err)
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
		exitFlag = true
	}

	if exitFlag {
		os.Exit(2)
	}

	log.Println("sACN-DMX is starting...")
	var devIndex int
	var devFound bool
	if CLIOptions.Device != "" {
		for i := range dl {
			if dl[i].SerialNumber == CLIOptions.Device {
				devIndex = i
				devFound = true
			}
		}
		if !devFound {
			log.Printf("Cannot find a device with S/N %s. Fallback to default...", CLIOptions.Device)
			devIndex = 0
		}
	} else {
		devIndex = 0
	}
	dmx, err := enttecOpenDMXConnect(dl[devIndex])
	if err == nil {
		log.Printf("Using %s (S/N %s)\n", dl[devIndex].Description, dl[devIndex].SerialNumber)
		defer dmx.Close()
	} else {
		log.Fatal(err)
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
		log.Println("Timeout detected on universe", univ)
		// Drop all DMX channels to zero on timeout
		if CLIOptions.Reset {
			for i := 0; i < 512; i++ {
				dmx.SetChannel((int16)(i+1), 0)
			}
		}
	})
	recv.Start()
	log.Printf("sACN listener started on %s.\n", CLIOptions.IPAddr)

	select {}
}
