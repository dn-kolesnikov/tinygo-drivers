// Package onewire provides a driver for 1-Wire devices over GPIO
package onewire

import (
	"errors"
	"machine"
	"time"
)

// OneWire ROM commands
const (
	ONEWIRE_READ_ROM   uint8 = 0x33
	ONEWIRE_MATCH_ROM  uint8 = 0x55
	ONEWIRE_SKIP_ROM   uint8 = 0xCC
	ONEWIRE_SEARCH_ROM uint8 = 0xF0
)

// Device wraps a connection to an 1-Wire devices.
type Device struct {
	Pin    machine.Pin
	NodeID []uint64
}

// Config wraps a configuration to an 1-Wire devices.
type Config struct{}

// Errors list
var (
	errNoPresence  = errors.New("Error: OneWire. No devices on the bus.")
	errReadAddress = errors.New("Error: OneWire. Read address error: CRC mismatch.")
)

// NewGPIO creates a new GPIO 1-Wire connection.
// The pin must be pulled up to the VCC via a resistor greater than 500 ohms (default 4.7k).
func New(pin machine.Pin, maxDevices uint8) Device {
	return Device{
		Pin:    pin,
		NodeID: make([]uint64, 0, maxDevices),
	}
}

// Configure initializes the protocol.
func (d *Device) Configure() {}

// Reset pull DQ line low, then up.
func (d Device) Reset() error {
	d.Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	time.Sleep(480 * time.Microsecond)
	d.Pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	time.Sleep(70 * time.Microsecond)
	precence := d.Pin.Get()
	time.Sleep(410 * time.Microsecond)
	if precence {
		return errNoPresence
	}
	return nil
}

// WriteBit transmits a bit to 1-Wire bus.
func (d Device) WriteBit(data uint8) {
	d.Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	if data&1 == 1 { // Send '1'
		time.Sleep(5 * time.Microsecond)
		d.Pin.Configure(machine.PinConfig{Mode: machine.PinInput})
		time.Sleep(60 * time.Microsecond)
	} else { // Send '0'
		time.Sleep(60 * time.Microsecond)
		d.Pin.Configure(machine.PinConfig{Mode: machine.PinInput})
		time.Sleep(5 * time.Microsecond)
	}
}

// Write transmits a byte as bit array to 1-Wire bus. (LSB first)
func (d Device) Write(data uint8) {
	for i := 0; i < 8; i++ {
		d.WriteBit(data)
		data >>= 1
	}
}

// ReadBit recieves a bit from 1-Wire bus.
func (d Device) ReadBit() (data uint8) {
	d.Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	time.Sleep(3 * time.Microsecond)
	d.Pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	time.Sleep(8 * time.Microsecond)
	if d.Pin.Get() {
		data = 1
	}
	time.Sleep(60 * time.Microsecond)
	return data
}

// Read recieves a byte from 1-Wire bus. (LSB first)
func (d Device) Read() (data uint8) {
	for i := 0; i < 8; i++ {
		data >>= 1
		data |= d.ReadBit() << 7
	}
	return data
}

// ReadAddress recieves a 64-bit unique ROM ID from Device. (LSB first)
// Note: use this if there is only one slave device on the bus.
func (d *Device) ReadAddress() error {
	d.NodeID = d.NodeID[:0]
	data := uint64(0)
	if err := d.Reset(); err != nil {
		return err
	}
	d.Write(ONEWIRE_READ_ROM)
	for i := 0; i < 64; i++ {
		data >>= 1
		data |= uint64(d.ReadBit()) << 63
	}
	if _, ok := Сrc8(data, 7); !ok {
		return errReadAddress
	}
	d.NodeID = append(d.NodeID, data)
	return nil
}

// Select selects the address of the device for communication
func (d Device) Select(index uint8) error {
	if err := d.Reset(); err != nil {
		return err
	}
	d.Write(ONEWIRE_MATCH_ROM)
	for i := 0; i < 64; i++ {
		d.WriteBit(uint8((d.NodeID[index] >> i) & 1))
	}
	return nil
}

// Search searches for all devices on the bus.
func (d *Device) Search() error {
	var (
		bit, bit_c  uint8
		bitOffset   uint8
		lastZero    uint8
		lastFork    uint8
		lastAddress uint64
	)
	if len(d.NodeID) > 0 {
		d.NodeID = d.NodeID[:0]
	}
	for ok := true; ok; ok = (lastFork != 0) {
		if err := d.Reset(); err != nil {
			return err
		}
		d.Write(ONEWIRE_SEARCH_ROM)
		for lastZero, bitOffset = 0, 0; bitOffset < 64; bitOffset++ {
			bit = d.ReadBit()           // read first address bit
			bit_c = d.ReadBit()         // read second (complementary) address bit
			if bit == 1 && bit_c == 1 { // no device
				return errNoPresence
			}
			if bit == 0 && bit_c == 0 { // collision
				if bitOffset == lastFork {
					bit = 1
				}
				if bitOffset < lastFork {
					bit = uint8((lastAddress >> uint64(bitOffset)) & 1)
				}
				if bit == 0 {
					lastZero = bitOffset
				}
			}
			if bit == 0 {
				lastAddress &= ^(1 << (uint64(bitOffset)))
			} else {
				lastAddress |= (1 << (uint64(bitOffset)))
			}
			d.WriteBit(bit)
		}
		d.NodeID = append(d.NodeID, lastAddress)
		lastFork = lastZero
	}
	return nil
}

// Crc8 compute a Dallas Semiconductor 8 bit CRC.
func Сrc8(buffer uint64, size int) (crc uint8, ok bool) {
	// Dow-CRC using polynomial X^8 + X^5 + X^4 + X^0
	// Tiny 2x16 entry CRC table created by Arjen Lentz
	// See http://lentz.com.au/blog/calculating-crc-with-a-tiny-32-entry-lookup-table
	crc8_table := [...]uint8{
		0x00, 0x5E, 0xBC, 0xE2, 0x61, 0x3F, 0xDD, 0x83,
		0xC2, 0x9C, 0x7E, 0x20, 0xA3, 0xFD, 0x1F, 0x41,
		0x00, 0x9D, 0x23, 0xBE, 0x46, 0xDB, 0x65, 0xF8,
		0x8C, 0x11, 0xAF, 0x32, 0xCA, 0x57, 0xE9, 0x74,
	}
	for i := 0; i < size; i++ {
		crc = uint8(buffer&0xFF) ^ crc // just re-using crc as intermediate
		crc = crc8_table[crc&0x0f] ^ crc8_table[16+((crc>>4)&0x0f)]
		buffer >>= 8
	}
	return crc, crc == uint8(buffer)
}
