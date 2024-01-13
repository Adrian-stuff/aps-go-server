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

	// Open GPIO PA_13 with pull-up resistor
	pin := allwinner.PA12
	errx := pin.In(gpio.PullDown, gpio.NoEdge)
	if errx != nil {
		fmt.Println("Error opening GPIO pin:", err)
		return
	}
	totalChan := make(chan int)
	go func() {
		total := 0
		value := 0
		counter := 0
		for {
			isReading := true

			for isReading {
				// pin.WaitForEdge(-1)

				pinState := pin.Read()

				if pinState == gpio.Low {
					counter++
					time.Sleep(100 * time.Millisecond) // Delay for 0.1 seconds
					// fmt.Println("counter: ", counter)

					if counter == 1 || counter == 3 || counter == 5 { // Check for specific counts immediately
						isReading = false

						// Calculate and print total based on counter value
						if counter == 1 || counter == 3 {
							value = 1
						} else { // Counter is 5
							value = 2
						}
						total += value
						// fmt.Println("total:", total)
						totalChan <- total
					}
				}
			}

			// Reset counter for the next cycle
			counter = 0
		}
	}()
	for {
		select {
		case <-totalChan:
			fmt.Printf("total: %d", totalChan)
		}
	}
}
