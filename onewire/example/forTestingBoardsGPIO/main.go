package main

import (
	"machine"
	"time"
)

func main() {

	// Define pin for test
	// the pin must be pulled up to the VCC via a resistor (default 4.7k).

	// for bluepill
	// pin := machine.PB12

	//for RP2040 pico
	pin := machine.GP16

	// for arduino
	// pin := machine.D12

	for {
		time.Sleep(3 * time.Second)

		println("Device:", machine.Device)

		start := time.Now()
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		end := time.Since(start)
		println("PinOutput:", pin.Get(), "Time in Microseconds:", end.Microseconds())

		start = time.Now()
		pin.Configure(machine.PinConfig{Mode: machine.PinInput})
		end = time.Since(start)
		println("PinInput:", pin.Get(), "Time in Microseconds:", end.Microseconds())

		println()
	}
}
