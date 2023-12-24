sACN to DMX converter
=====================

This little command-line tool provides an interface between sACN network and DMX hardware. Basically it feeds your DMX dongle with whatever comes from sACN receiver.

At the moment being there are two supported types of interfaces:
1. Enttec Open DMX. It is widely known as one of the cheapest DMX dongles on the market. Many clones of this device have been manufactured around the world thanks to its simplicity and openness. Basically, a variety of cheap USB-RS485 converters off aliexpress will work just as good as original device (as soon as it is a FTDI chip based device). You can even make one yourself if you know which side of soldering iron is safe to touch :) I will not provide instructions though. Google it.

2. Generic UART interface. As soon as your UART (COM-port) interface can run at 250kbit it can be used (in theory) to feed DMX data to fixtures. Not all built-in UARTs are able to run that fast, but it's not a good idea to wire your motherboard directly to your lights anyway. Yet i tested this approach on FT232 VCP mode, CH340/341 USB-UART chips and built-in UART on Raspberry Pi and it works. So if previous option is not for you, you can test this one. It is supposed to work with pretty much any USB-RS485 interface. No warranties provided though, i only tested it on what i have on hand.

How to use
----------

This is a command-line tool, so you run this in terminal. If you're a linux user you need no explaination. If you run Windows search for "command line" or "power shell" in your Start menu.

First you want to find out what hardware devices are available in your system. You may run the tool with `-t <device type> -l` options:

	d:\sacndmx\sacndmx.exe -t uart -l
	Port 0: COM1
	Port 1: COM6

Now you see what devices are available. There are two possible variants of '-t': uart and opendmx. Latter is default if no `-t` option provided.

Also you may want to know what ip-addresses you can bind to. Check possible adresses with `-n` option:

	d:\sacndmx\sacndmx.exe -n
	Local addr: 169.254.135.160
	Local addr: 169.254.244.130
	Local addr: 192.168.1.101

Now let's say you want to run the tool on COM6 and bind to 192.168.1.101. Provide `-t`, `-d` and `-s` options:

	D:\sacndmx>sacndmx.exe -t uart -d COM6 -s 192.168.1.101
	sACN-DMX is starting...
	Using general UART (COM6)
	sACN listener started on 192.168.1.101.
	DMX stream started on general UART (COM6)

In most cases you will need to provide an ip address option. By default the tool will stick itself to 'localhost' address, which means that it will only listen to a machine it is running on.

On linux, you use your tty device path instaed of port name:

	lightmaster@pi ~/sacndmx> ./sacndmx -t uart -d /dev/ttyAMA0 -s 192.168.1.90
	sACN-DMX is starting...
	Using general UART (/dev/ttyAMA0)
	sACN listener started on 192.168.1.90.
	DMX stream started on general UART (/dev/ttyAMA0)

And finally, if you want to run Enttec OpenDMX'ish device, you can let the `-t` option go and you can pick a particular device by its serial number:

	lightmaster@pi ~/sacndmx> ./sacndmx -l
    Device 0: Uncle Eugene's DMX Dongle A (S/N A5069FNKA)
    Device 1: Uncle Eugene's DMX Dongle B (S/N A5069FNKB)
    Device 2: Uncle Eugene's DMX Dongle C (S/N A5069FNKC)
    Device 3: Uncle Eugene's DMX Dongle D (S/N A5069FNKD)

    lightmaster@pi ~/sacndmx> ./sacndmx -s 192.168.11.100
    sACN-DMX is starting...
    Using Uncle Eugene's DMX Dongle A (A5069FNKA)
    sACN listener started on 192.168.1.101.
    DMX stream started on Uncle Eugene's DMX Dongle A (A5069FNKA)

All this is just the same way on Windows. You can of course get full list of possible options with usual `-h` key.

Not yet here
------------

There are couple of things not yet implemented but planned.
1. incoming Art-Net support. Must be an easy one, there are libraries :)
2. More hardware types. uDMX and enttecPro to begin with.
3. Automatic hardware reconnect. Doable, but will take some effort to implement properly. Right now it will just fail if hardware connection is lost.
4. Web interface. Well, maybe one day, idk :)

Licensing and warranties
------------------------

This tool is licensed on hell knows what kind of license. Probably some GPL as i use third-party libraries. So i assume you're free to download and modify it for your own sake. And i provide absolutely no warranties. It may work or it may break right in the middle of important show. We just don't know there's not enought statistical data :) Although i've tested the tool myself and actually use it, you better give it a good test run on your hardware before any production use.
