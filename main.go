package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/philippebeaulieu/rpi-thermostat/api"
	"github.com/philippebeaulieu/rpi-thermostat/controller/rpi"
	"github.com/philippebeaulieu/rpi-thermostat/database"
	"github.com/philippebeaulieu/rpi-thermostat/datagatherer"
	"github.com/philippebeaulieu/rpi-thermostat/sensor/ds18b20"
	"github.com/philippebeaulieu/rpi-thermostat/thermostat"
	"github.com/philippebeaulieu/rpi-thermostat/weather/apixu"
)

func main() {
	sensor, err := ds18b20.NewDs18b20("28-041685fc45ff")
	// sensor, err := fakesensor.NewFakeSensor()
	if err != nil {
		fmt.Printf("failed to create sensor: %v\n", err)
		return
	}

	controller, err := rpi.NewRpiController()
	// controller, err := fakecontroller.NewFakeController()
	if err != nil {
		fmt.Printf("failed to create controller: %v\n", err)
		return
	}

	thermostat := thermostat.NewThermostat(sensor, controller)
	go thermostat.Run()

	weather := apixu.NewApixuWeather(thermostat, "61a17a8fdb264c2eaba152957173006", "J7J0B7")
	go weather.Run()

	database, err := database.NewDatabase(thermostat, "192.168.2.41:3306", "thermostat", "GDeWFE8Hg3aKh44")
	if err != nil {
		fmt.Printf("failed to create database: %v\n", err)
		return
	}

	datagatherer, err := datagatherer.NewDataGatherer(thermostat, database)
	if err != nil {
		fmt.Printf("failed to create datagatherer: %v\n", err)
		return
	}

	go datagatherer.Run()

	apiserver := apiserver.NewAPIServer(thermostat, datagatherer, 80)
	// apiserver := apiserver.NewAPIServer(thermostat, datagatherer, 8080)
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
	controller.Off(1)
	controller.Off(2)
	controller.Off(3)
}
