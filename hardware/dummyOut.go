// This file implements dummy out interface. It dumps output to console
// and it can be used for testing purposes when do device is available

package hardware

import (
	"fmt"
	"strconv"
)

var helicopter = []byte{'|', '/', '-', '\\'}
var heliphase byte

// EnttecOpenDMX is a structure representing openDMX entity
type dummyOut struct {
	mode    string
	channel []byte
	buffer  []byte
}

// method Init looks up for available hardware devices
func DummyOutInit() (*dummyOut, error) {
	d := new(dummyOut)
	return d, nil
}

func (d *dummyOut) List() error {
	fmt.Println("Dummy output.")
	fmt.Println("There are two possible options:")
	fmt.Println("1. You can specify 'ti' device. 'ti' stcands for 'traffic indicator'.")
	fmt.Println("   In this mode you will get simple rotating traffic indicator in console.")
	fmt.Println("2. You can specify a number of channels to show. This one is default.")
	fmt.Println("   Any non number input will lead to all 512 channels to dump.")
	return nil
}

func (d *dummyOut) GetSerial() string {
	return d.mode
}

func (d *dummyOut) GetDescription() string {
	return "console"
}

func (d *dummyOut) SelectDevice(mode string) error {
	d.mode = mode
	return nil
}

func (d *dummyOut) Connect() error {
	return nil
}

// method Close closes a device connection
func (d *dummyOut) Close() error {
	return nil
}

// method SetChannel sets given DMX channel a value
func (d *dummyOut) SetChannel(index int16, data byte) error {
	var num int16 = 512
	if d.mode == "ti" {
		fmt.Printf("\r%c", helicopter[heliphase])
		heliphase++
		if heliphase == 4 {
			heliphase = 0
		}
	} else {
		if n, err := strconv.Atoi(d.mode); err == nil {
			num = (int16)(n)
		}
		if index <= num {
			fmt.Printf("DMX Data: %3d - %3d\n", index, data)
		}
	}
	return nil
}

// Method run is intended to start as a goroutine and takes care of device communication
func (d *dummyOut) Run() {
	select {}
}

// private methods section

// method Render processes the DMX magic
func (d *dummyOut) render() error {
	return nil
}
