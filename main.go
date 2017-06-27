package main

import (
	"fmt"
	"os"
	"os/signal"

	"./api"
	"./controller/rpi"
	"./sensor/ds18b20"
	"./thermostat"
)

func main() {
	sensor, err := ds18b20.NewDs18b20()
	if err != nil {
		fmt.Printf("failed to create sensor: %v\n", err)
		return
	}

	controller, err := rpi.NewRpiController()
	if err != nil {
		fmt.Printf("failed to create controller: %v\n", err)
		return
	}

	thermostat, err := thermostat.NewThermostat(sensor, controller, 21)
	if err != nil {
		fmt.Printf("failed to create controller: %v\n", err)
		return
	}
	go thermostat.Run()

	apiserver := apiserver.NewAPIServer(thermostat)
	go apiserver.Run()

	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		close(done)
	}()
	<-done
	fmt.Printf("shutting down")
	controller.Off()
}
