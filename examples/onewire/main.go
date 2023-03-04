package main

import (
	"encoding/hex"
	"machine"
	"time"

	"github.com/dn-kolesnikov/tinygo-drivers/onewire"
)

func main() {
	// Define pin for 1-wire devices
	// the pin must be pulled up to the VCC via a resistor (default 4.7k).

	// for bluepill
	// pin := machine.PA0

	//for RP2040 pico
	pin := machine.GP16

	// for arduino
	// pin := machine.PB2

	ow := onewire.New(pin)

	for {
		time.Sleep(3 * time.Second)

		println()
		println("Device:", machine.Device)

		romIDs, err := ow.Search(onewire.ONEWIRE_SEARCH_ROM)
		if err != nil {
			println(err)
		}
		for _, romid := range romIDs {
			println(hex.EncodeToString(romid))
		}

		if len(romIDs) == 1 {
			// only 1 device on bus
			r, err := ow.ReadAddress()
			if err != nil {
				println(err)
			}
			println(hex.EncodeToString(r))

		}

	}
}
