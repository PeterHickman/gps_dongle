I bought a GPS dongle to play with which required to talk to the dongle as a USB modem. The C code here is modified from the product [wiki](https://wiki.dfrobot.com/USB_GPS_Receiver_SKU_TEL0137)

The go version was hacked up for the hell of it

Item of note. On the Macintosh the dongle appeared as both `/dev/cu.usbmodem1101` and `/dev/tty.usbmodem1101`. But the `/dev/cu.*` device was the only one that worked
