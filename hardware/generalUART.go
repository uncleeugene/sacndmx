package hardware

import (
	"errors"
	"fmt"
	"time"

	"go.bug.st/serial"
)

var uartMode = serial.Mode{
	BaudRate: 250000,
	DataBits: 8,
	Parity:   serial.NoParity,
	StopBits: 2,
}

type UART struct {
	selected int
	ports    []string
	port     serial.Port
	channel  []byte
	buffer   []byte
}

func UartInit() (*UART, error) {
	u := new(UART)
	if ports, err := serial.GetPortsList(); err == nil {
		u.ports = ports
		return u, nil
	} else {
		return nil, err
	}
}

func (d *UART) List() error {
	if len(d.ports) != 0 {
		for i := range d.ports {
			fmt.Printf("Port %d: %s\n", i, d.ports[i])
		}
		return nil
	} else {
		return errors.New("no ports found")
	}

}

func (d *UART) GetSerial() string {
	return d.ports[d.selected]
}

func (d *UART) GetDescription() string {
	return "general UART"
}

func (d *UART) SelectDevice(name string) error {
	for i := range d.ports {
		if d.ports[i] == name {
			d.selected = i
			return nil
		}
	}
	return errors.New("specified device not found")
}

func (d *UART) Connect() error {
	// Connect method establishes a connection to a device
	if len(d.ports) != 0 {
		port, err := serial.Open(d.ports[d.selected], &uartMode)
		if err != nil {
			return err
		}
		d.port = port
		d.channel = make([]byte, 512)
		d.buffer = make([]byte, 513)
		return nil
	} else {
		return errors.New("no devices found")
	}
}

// method Close closes a device connection
func (d *UART) Close() error {
	return d.port.Close()
}

// method SetChannel sets given DMX channel a value
func (d *UART) SetChannel(index int16, data byte) error {
	d.channel[index-1] = data
	return nil
}

// Method run is intended to start as a goroutine and takes care of device communication
func (d *UART) Run() {
	for {
		if err := d.render(); err != nil {
			fmt.Printf("Error rendering DMX: %s\n", err)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// private methods section

// method Render processes the DMX magic
func (d *UART) render() error {
	for i := 0; i < 512; i++ {
		d.buffer[i+1] = d.channel[i]
	}
	d.setBreak()
	if _, err := d.port.Write(d.buffer); err != nil {
		return err
	}
	return nil
}

func (d *UART) setBreak() {
	var breakMode = serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: 0,
	}
	d.port.SetMode(&breakMode)
	d.port.Write([]byte("0"))
	d.port.SetMode(&uartMode)
}
