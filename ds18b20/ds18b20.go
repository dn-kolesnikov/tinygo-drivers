package ds18b20

import (
	"errors"

	"github.com/dn-kolesnikov/tinygo-drivers/onewire"
)

// Device ROM commands
const (
	DS18B20_CONVERT_TEMPERATURE uint8 = 0x44
	DS18B20_READ_SCRATCHPAD     uint8 = 0xBE
	DS18B20_COPY_SCRATCHPAD     uint8 = 0x48
	DS18B20_WRITE_SCRATCHPAD    uint8 = 0x4E
	DS18B20_READ_POWER_SUPPLY   uint8 = 0xB4
	DS18B20_RECALL_E2           uint8 = 0xB8
)

// Device wraps a connection to an 1-Wire devices.
type Device struct {
	owd onewire.Device
}

// Errors list
var (
	errReadTemperature = errors.New("Error: DS18B20. Read temperature error: CRC mismatch.")
)

//
func New(owd onewire.Device) Device {
	return Device{
		owd: owd,
	}
}

// Configure. Initializes the device, left for compatibility reasons.
func (d Device) Configure() {}

// ThermometerResolution sets thermometer resolution from 9 to 12 bits
func (d Device) ThermometerResolution(romid []uint8, resolution uint8) {
	if 9 <= resolution && resolution <= 12 {
		d.owd.Select(romid)
		d.owd.Write(DS18B20_WRITE_SCRATCHPAD)       // send three data bytes to scratchpad (TH, TL, and config)
		d.owd.Write(0xFF)                           // to TH
		d.owd.Write(0x00)                           // to TL
		d.owd.Write(((resolution - 9) << 5) | 0x1F) // to resolution config
	}
}

// RequestTemperature sends request to device
func (d Device) RequestTemperature(romid []uint8) {
	d.owd.Select(romid)
	d.owd.Write(DS18B20_CONVERT_TEMPERATURE)
}

// ReadTemperatureRaw returns the raw temperature.
// ScratchPad memory map:
// byte 0: Temperature LSB
// byte 1: Temperature MSB
func (d Device) ReadTemperatureRaw(romid []uint8) ([]uint8, error) {
	var spb = make([]uint8, 9) // ScratchPad buffer
	d.owd.Select(romid)
	d.owd.Write(DS18B20_READ_SCRATCHPAD)
	for i := 0; i < 9; i++ {
		spb[i] = d.owd.Read()
	}
	if onewire.Сrc8(spb, 8) != spb[8] {
		return []uint8{}, errReadTemperature
	}
	return spb[:2], nil
}

// ReadTemperature returns the temperature in celsius milli degrees (°C/1000)
func (d Device) ReadTemperature(romid []uint8) (int32, error) {
	raw, err := d.ReadTemperatureRaw(romid)
	if err != nil {
		return 0, err
	}
	t := int32(uint16(raw[0]) | uint16(raw[1])<<8)
	if t&0x8000 == 0x8000 {
		t -= 0x10000
	}
	return (t * 625 / 10), nil
}
