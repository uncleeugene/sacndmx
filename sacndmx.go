package main

import (
	"flag"
	"log"

	"github.com/oliread/usbdmx"
	"github.com/oliread/usbdmx/ft232"
)

func main() {

	vid := uint16(0x0403)
	pid := uint16(0x6001)
	inputInterfaceID := flag.Int("input-id", 0, "Input interface ID for device")
	outputInterfaceID := flag.Int("output-id", 0, "Output interface ID for device")
	debugLevel := flag.Int("debug", 0, "Debug level for USB context")
	flag.Parse()

	// Create a configuration from our flags
	config := usbdmx.NewConfig(vid, pid, *inputInterfaceID, *outputInterfaceID, *debugLevel)

	// Get a usb context for our configuration
	config.GetUSBContext()

	// Create a controller and connect to it
	controller := ft232.NewDMXController(config)
	if err := controller.Connect(); err != nil {
		log.Fatalf("Failed to connect DMX Controller: %s", err)
	} else {
		log.Printf("Connected to OpenDMX!")
	}
}
