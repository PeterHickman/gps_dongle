package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PeterHickman/toolbox"
	"go.bug.st/serial"
)

var list bool
var device string

func init() {
	l := flag.Bool("list", false, "List the available usb modems")
	d := flag.String("device", "", "Read from this device")

	flag.Parse()

	list = *l
	device = *d

	if !list && !toolbox.FileExists(device) {
		log.Fatal("Cannot find the supplied device")
	}
}

func list_ports() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	fmt.Println("List of available USB modems")
	for _, port := range ports {
		fmt.Printf("  %v\n", port)
	}
}

func split_line(line string) []string {
	l := line

	i := strings.Index(l, "*")
	if i != -1 {
		l = line[:i]
	}

	return strings.Split(l, ",")
}

func decode_gga(line string) {
	parts := split_line(line)

	// UTC Time (hhmmss.ss)
	// Latitude (ddmm.mmmm)
	// N/S Indication
	// Longitude (ddmm.mmmm)
	// E/W Indication

	fmt.Printf("GGA %s %s %s %s %s\n", parts[1], parts[2], parts[3], parts[4], parts[5])
}

func decode_rmc(line string) {
	parts := split_line(line)

	// UTC Time (hhmmss.ss)
	// Latitude (ddmm.mmmm)
	// N/S Indication
	// Longitude (dddmm.mmmm)
	// E/W Indication

	fmt.Printf("RMC %s %s %s %s %s\n", parts[1], parts[3], parts[4], parts[5], parts[6])
}

func decode_message(block string) {
	parts := strings.Split(block, "\n")
	for _, line := range parts {
		if strings.HasPrefix(line, "$") {
			switch line[3:6] {
			case "GGA":
				decode_gga(line)
			case "RMC":
				decode_rmc(line)
			}
		}
	}
}

func main() {
	if list {
		list_ports()
		os.Exit(0)
	}

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(device, mode)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 600)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		decode_message(string(buff[:n]))
	}

	fmt.Println("Done")
}
