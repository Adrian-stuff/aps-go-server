package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/allwinner"
)

func monitorPullout(pin gpio.PinIO, trigger chan<- gpio.Level) {
	for {
		pin.WaitForEdge(-1)   // Wait for falling edge on the pin
		trigger <- pin.Read() // Send trigger signal to channel
	}
}
func main() {
	// Initialize the periph library
	fmt.Println("run")

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open GPIO PA_12 with pull-up resistor
	pin := allwinner.PA13
	err := pin.In(gpio.PullUp, gpio.FallingEdge) // Note the use of orangepi.GPIO12
	if err != nil {
		fmt.Println("Error opening GPIO pin:", err)
		return
	}
	// Create channel for pullout triggers
	triggerChan := make(chan gpio.Level)

	// Launch goroutine for edge detection
	go monitorPullout(pin, triggerChan)

	// Main program loop
	for {
		select {
		case <-triggerChan:
			fmt.Printf("Pullout triggered! %t \n", triggerChan)
			// Handle the pullout event here
		default:
			// Perform other tasks while waiting for trigger
			time.Sleep(time.Second / 2) // Adjust sleep time as needed
		}
	}
}
