package main

import (
	"machine"
	"time"

	"github.com/dn-kolesnikov/tinygo-drivers/ds18b20"
	"github.com/dn-kolesnikov/tinygo-drivers/onewire"
)

func main() {
	// Define pin for DS18B20
	// the pin must be pulled up to the VCC via a resistor (default 4.7k).

	// for bluepill
	// pin := machine.PA0

	//for RP2040 pico
	pin := machine.GP16

	// for arduino
	// pin := machine.D12

	ow := onewire.New(pin)
	sensor := ds18b20.New(ow)

	for {
		time.Sleep(3 * time.Second)

		println()
		println("Device:", machine.Device)

		println("Request Temperature.")
		err := sensor.RequestTemperature()
		if err != nil {
			println(err)
		}

		// wait 750ms or more for DS18B20 convert T
		time.Sleep(1 * time.Second)

		println("Read Temperature")
		t, err := sensor.ReadTemperature()
		if err != nil {
			println(err)
		}
		// temperature in celsius milli degrees (°C/1000)
		println("Temperature (°C/1000): ", t)
	}
}
