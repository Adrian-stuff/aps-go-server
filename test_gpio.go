package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/allwinner"
)

func main() {
	// Initialize periph.io with Orange Pi One support
	_, err := host.Init()
	if err != nil {
		fmt.Println("Error initializing periph.io:", err)
		return
	}

	// Ensure Orange Pi One support is loaded
	// if _, err := orangepi.Present(); err != nil {
	// 	fmt.Println("Error: Orange Pi One not supported by periph.io")
	// 	return
	// }

	// Open GPIO PA_13 with pull-up resistor
	pin := allwinner.PA13
	pin.In(gpio.PullUp, gpio.BothEdges)
	// if err != nil {
	// 	fmt.Println("Error opening GPIO pin:", err)
	// 	return
	// }

	// Main loop
	for {
		isReading := true
		counterTotal := 0
		total := 0
		value := 0
		counter := 0

		for isReading {
			pinState := pin.Read()

			if pinState == gpio.Low {
				counter++
				time.Sleep(100 * time.Millisecond) // Delay for 0.1 seconds
				fmt.Println("counter: ", counter)

				counterTotal = counter
				if counterTotal == 1 || counterTotal == 3 || counterTotal == 5 {
					isReading = false
				}
			}
		}

		if counterTotal == 1 || counterTotal == 3 {
			value++ // Increment value directly (equivalent to value + 1)
		} else if counterTotal == 5 {
			value += 2 // Add 2 for counterTotal == 5
		}

		total += value

		if total != 0 {
			fmt.Println("total:", total)
		}

		// Reset variables for the next cycle
		total = 0
		value = 0
		counter = 0
	}
}
