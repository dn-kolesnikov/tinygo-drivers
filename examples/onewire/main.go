package main

import (
	"encoding/hex"
	"machine"
	"time"

	"github.com/dn-kolesnikov/tinygo-drivers/onewire"
)

func main() {

	//for RP2040 pico
	pin := machine.GP16

	ow := onewire.New(pin)

	for {
		time.Sleep(3 * time.Second)

		println()
		println("Device:", machine.Device)

		romIDs, err := ow.Search(onewire.SEARCH_ROM_COMMAND)
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
