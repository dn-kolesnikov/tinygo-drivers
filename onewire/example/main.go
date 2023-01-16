package main

import (
	"machine"
	"time"

	"github.com/dn-kolesnikov/onewire"
	"github.com/dn-kolesnikov/onewire/devices/ds18b20"
)

func main() {
	// Define pin for DS18B20
	// the pin must be pulled up to the VCC via a resistor (default 4.7k).

	// for bluepill
	// pin := machine.PB12

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

		println("Read OneWire ROM.")
		println("Send command =", SliceToHexString([]uint8{onewire.ONEWIRE_READ_ROM}))
		err := sensor.ReadAddress()
		if err != nil {
			println(err)
			continue
		}
		println("RomID:", SliceToHexString(sensor.RomID))

		println("Request Temperature.")
		println("Send command =", SliceToHexString([]uint8{ds18b20.DS18B20_CONVERT_TEMPERATURE}))
		err = sensor.RequestTemperature()
		if err != nil {
			println(err)
			continue
		}

		// wait 750ms or more for DS18B20 convert T
		time.Sleep(1 * time.Second)

		println("Read Temperature")
		println("Send command =", SliceToHexString([]uint8{ds18b20.DS18B20_READ_SCRATCHPAD}))
		t, err := sensor.ReadTemperature()
		if err != nil {
			println(err)
			continue
		}
		// temperature in celsius milli degrees (°C/1000)
		println("Temperature (°C/1000): ", t)
	}
}

// Convert a slice to Hex string
// fmt.Printf - compile error on an Arduino Uno boards
func SliceToHexString(rom []uint8) string {
	const hc string = "0123456789ABCDEF"
	var result string = "0x"
	for _, v := range rom {
		if v < 0x10 {
			result += "0" + string(hc[v])
		} else {
			result += string(hc[v&0xF0>>4]) + string(hc[v&0x0F])
		}
	}
	return result
}
