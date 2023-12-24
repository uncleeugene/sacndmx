package hardware

type Hardware interface {
	Connect() error
	Close() error
	SetChannel(int16, byte) error
	Run()
	List() error
	GetDescription() string
	GetSerial() string
	SelectDevice(string) error
}
