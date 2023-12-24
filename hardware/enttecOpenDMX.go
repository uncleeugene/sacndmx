// This file implements Enttec Open DMX interface

package hardware

import (
	"errors"
	"fmt"
	"time"

	ftdi "github.com/uncleeugene/goftdi"
)

// Line properties are hardcoded because DMX line pearameters are hardcoded and won't change anyways
var lineProperties = ftdi.LineProperties{
	Bits:     ftdi.BITS_8,
	StopBits: ftdi.STOP_2,
	Parity:   ftdi.NONE,
}

// EnttecOpenDMX is a structure representing openDMX entity
type EnttecOpenDMX struct {
	selected int
	devList  []ftdi.DeviceInfo
	device   *ftdi.Device
	channel  []byte
	buffer   []byte
}

// method Init looks up for available hardware devices
func EnttecOpenDMXInit() (*EnttecOpenDMX, error) {
	d := new(EnttecOpenDMX)
	dl, err := ftdi.GetDeviceList()
	d.devList = dl
	return d, err
}

func (d *EnttecOpenDMX) List() error {
	if len(d.devList) != 0 {
		for i := range d.devList {
			fmt.Printf("Device %d: %s (S/N %s)\n", i, d.devList[i].Description, d.devList[i].SerialNumber)
		}
		return nil
	} else {
		return errors.New("no devices found")
	}

}

func (d *EnttecOpenDMX) GetSerial() string {
	return d.devList[d.selected].SerialNumber
}

func (d *EnttecOpenDMX) GetDescription() string {
	return d.devList[d.selected].Description
}

func (d *EnttecOpenDMX) SelectDevice(serial string) error {
	for i := range d.devList {
		if d.devList[i].SerialNumber == serial {
			d.selected = i
			return nil
		}
	}
	return errors.New("OpenDMX: no device with selected S/N found")
}

func (d *EnttecOpenDMX) Connect() error {
	// Connect method establishes a connection to a device
	if len(d.devList) != 0 {
		dev, err := ftdi.Open(d.devList[d.selected])
		if err != nil {
			return err
		}
		d.device = dev
		d.channel = make([]byte, 512)
		d.buffer = make([]byte, 513)
		d.device.SetLineProperty(lineProperties)
		d.device.SetBaudRate(250000)
		d.device.Purge()
		return nil
	} else {
		return errors.New("no devices found")
	}
}

// method Close closes a device connection
func (d *EnttecOpenDMX) Close() error {
	return d.device.Close()
}

// method SetChannel sets given DMX channel a value
func (d *EnttecOpenDMX) SetChannel(index int16, data byte) error {
	d.channel[index-1] = data
	return nil
}

// Method run is intended to start as a goroutine and takes care of device communication
func (d *EnttecOpenDMX) Run() {
	for {
		if err := d.render(); err != nil {
			fmt.Printf("Error rendering DMX: %s\n", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// private methods section

// method Render processes the DMX magic
func (d *EnttecOpenDMX) render() error {
	for i := 0; i < 512; i++ {
		d.buffer[i+1] = d.channel[i]
	}
	d.device.SetBreakOn(lineProperties)
	d.device.SetBreakOff(lineProperties)
	if _, err := d.device.Write(d.buffer); err != nil {
		return err
	}
	return nil
}
