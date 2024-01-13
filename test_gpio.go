package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/orangepi"
)

func monitorPullout(pin gpio.PinIO, trigger chan<- bool) {
	for {
		pin.WaitForEdge(-1) // Wait for falling edge on the pin
		trigger <- true     // Send trigger signal to channel
	}
}
func main() {
	// Initialize the periph library
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open GPIO PA_12 with pull-up resistor
	pin := orangepi.PA1_12
	err := pin.In(gpio.PullUp, gpio.FallingEdge) // Note the use of orangepi.GPIO12
	if err != nil {
		fmt.Println("Error opening GPIO pin:", err)
		return
	}
	// Create channel for pullout triggers
	triggerChan := make(chan bool)

	// Launch goroutine for edge detection
	go monitorPullout(pin, triggerChan)

	// Main program loop
	for {
		select {
		case <-triggerChan:
			fmt.Println("Pullout triggered!")
			// Handle the pullout event here
		default:
			// Perform other tasks while waiting for trigger
			time.Sleep(time.Second) // Adjust sleep time as needed
		}
	}
}
