# 1-Wire Driver for DS18B20 Digital Thermometer in TinyGo

This repository provides a **1-Wire driver** implementation for the **DS18B20 digital thermometer**, specifically designed for use with **TinyGo**. The DS18B20 is a popular temperature sensor that communicates over the 1-Wire protocol, making it ideal for low-power and embedded applications.

The driver has been tested on several development boards, ensuring compatibility and performance across different hardware platforms.

---

## About 1-Wire Protocol

The **1-Wire protocol** is a single-wire communication interface developed by Maxim Integrated (formerly Dallas Semiconductor). It allows devices to communicate over a single data line, making it simple and cost-effective for connecting sensors like the DS18B20. Learn more about 1-Wire on [Wikipedia](https://en.wikipedia.org/wiki/1-Wire).

---

## Features

- **DS18B20 Support**: Fully compatible with the DS18B20 digital thermometer.
- **TinyGo Compatibility**: Designed to work seamlessly with TinyGo, enabling Go programming on microcontrollers.
- **Cross-Platform**: Tested on multiple development boards (see below for details).
- **Efficient Communication**: Implements the 1-Wire protocol efficiently, ensuring reliable data transfer.
- **Easy Integration**: Simple API for reading temperature data from the DS18B20 sensor.

---

## Tested Development Boards

The driver has been tested on the following boards:

- **Raspberry Pi Pico (RP2040)** - **Working**
- **WeAct RP2040 Pico** - **Working**
- **Seeed Studio XIAO RP2040** - **Working**
- **Arduino UNO** - **Working**  
  **Note**: Some cheap Arduino UNO clones may not work due to slow GPIO level switching, which can cause freezing.
- **BluePill (STM32F103C8T6)** - **Not Working**  
  **Issue**: GPIO levels on the 1-Wire bus do not switch correctly. This may require additional GPIO configuration or hardware-specific adjustments.

---

## Getting Started

### Prerequisites

- **TinyGo**: Ensure you have TinyGo installed. Follow the [TinyGo installation guide](https://tinygo.org/getting-started/).
- **DS18B20 Sensor**: Connect the DS18B20 sensor to your development board using the 1-Wire protocol.
- **Pull-Up Resistor**: A 4.7kΩ pull-up resistor is required on the 1-Wire data line for reliable communication.

### Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/dn-kolesnikov/tinygo-drivers.git
   ```
2. Navigate to the repository directory:
   ```bash
   cd tinygo-drivers
   ```
3. Build and flash the example to your board:
   ```bash
   tinygo flash -target=<your-board> examples/ds18b20
   ```
## Example Code

Here’s a simple example to read temperature data from the DS18B20 sensor:

```go
ackage main

import (
	"encoding/hex"
	"machine"
	"time"

	"github.com/dn-kolesnikov/tinygo-drivers/ds18b20"
	"github.com/dn-kolesnikov/tinygo-drivers/onewire"
)

func main() {
	// Define pin for DS18B20

	//for RP2040 pico
	pin := machine.GP16

	ow := onewire.New(pin)
	romIDs, err := ow.Search(onewire.SEARCH_ROM_COMMAND)
	if err != nil {
		println(err)
	}
	sensor := ds18b20.New(ow)

	for {
		time.Sleep(3 * time.Second)

		println()
		println("Device:", machine.Device)

		println()
		println("Request Temperature.")
		for _, romid := range romIDs {
			println("Sensor RomID: ", hex.EncodeToString(romid))
			sensor.RequestTemperature(romid)
		}

		// wait 750ms or more for DS18B20 convert T
		time.Sleep(1 * time.Second)

		println()
		println("Read Temperature")
		for _, romid := range romIDs {
			raw, err := sensor.ReadTemperatureRaw(romid)
			if err != nil {
				println(err)
			}
			println()
			println("Sensor RomID: ", hex.EncodeToString(romid))
			println("Temperature Raw value: ", hex.EncodeToString(raw))

			t, err := sensor.ReadTemperature(romid)
			if err != nil {
				println(err)
			}
			println("Temperature in celsius milli degrees (°C/1000): ", t)
		}
	}
}
```

## Known Issues and Limitations

- **BluePill Compatibility**: The driver does not currently work on the BluePill board due to GPIO level switching issues. This may require further investigation or hardware-specific adjustments.
- **Cheap Arduino UNO Clones**: Some low-cost Arduino UNO clones may not work reliably due to slow GPIO switching, which can cause the program to freeze.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request. Ensure your code follows the project’s coding standards and includes appropriate tests.

## License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute it as needed.

## Acknowledgments

- **TinyGo Community**: For creating an amazing Go compiler for microcontrollers.
- **Maxim Integrated**: For developing the 1-Wire protocol and DS18B20 sensor.
- **Testers**: Thanks to everyone who tested the driver on various boards and provided feedback.

## Resources

- [1-Wire Protocol (Wikipedia)](https://en.wikipedia.org/wiki/1-Wire)
- [DS18B20 Datasheet](https://datasheets.maximintegrated.com/en/ds/DS18B20.pdf)
- [TinyGo Documentation](https://tinygo.org/docs/)
- [Raspberry Pi Pico Documentation](https://www.raspberrypi.com/documentation/microcontrollers/)
