package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/philippebeaulieu/rpi-thermostat/api"
	"github.com/philippebeaulieu/rpi-thermostat/controller/rpi"
	"github.com/philippebeaulieu/rpi-thermostat/database"
	"github.com/philippebeaulieu/rpi-thermostat/sensor/ds18b20"
	"github.com/philippebeaulieu/rpi-thermostat/thermostat"
)

func main() {
	sensor, err := ds18b20.NewDs18b20("28-041685fc45ff")
	if err != nil {
		fmt.Printf("failed to create sensor: %v\n", err)
		return
	}

	controller, err := rpi.NewRpiController()
	if err != nil {
		fmt.Printf("failed to create controller: %v\n", err)
		return
	}

	thermostat := thermostat.NewThermostat(sensor, controller, 21)
	go thermostat.Run()

	apiserver := apiserver.NewAPIServer(thermostat, 80)
	go apiserver.Run()

	database, err := database.NewDatabase(thermostat)
	go database.Run()

	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		close(done)
	}()
	<-done
	fmt.Printf("shutting down")
	controller.Off(1)
	controller.Off(2)
	controller.Off(3)
}
