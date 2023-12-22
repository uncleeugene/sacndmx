This is a small tool that listens to a network for sACN 1.31 traffic and sends all it gets to Enttec Open DMX device.

Runs on Windows and linux. arm6 binary is included, it is confirmed to run on Raspberry Pi (even old Model A one). I don't see why it can't run on Mac, but i have no Mac on hands to test.

As soon as Enttec Open DMX is FTDI FT232 based device software utilizes D2XX drivers on Windows and on linux libftdi is used.

This is a very homebrew tool as of now, i made if just for my own purposes, and i cannot guarantee it will run steady enough for production use. Although it runs very smooth in my home pool.

There's no documentation so far, only -h flag and it is a bit ahead of time. There's no config file :)

You could list your network addresses by using -i flag, and you can list your FTDI devises using -f flag.

Run the tool with sacndmx -a \<binding ip\>

If you need to specify a device to output to, you have -d \<serial number\> flag.

-r flag will force the software to shut down all DMX channels in case there's timeout in sACN. By default DMX values will hold.

That's pretty much it :)
