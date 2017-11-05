package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/philippebeaulieu/rpi-thermostat/api"
	"github.com/philippebeaulieu/rpi-thermostat/controller"
	"github.com/philippebeaulieu/rpi-thermostat/controller/fakecontroller"
	"github.com/philippebeaulieu/rpi-thermostat/controller/rpi"
	"github.com/philippebeaulieu/rpi-thermostat/database"
	"github.com/philippebeaulieu/rpi-thermostat/datagatherer"
	"github.com/philippebeaulieu/rpi-thermostat/sensor"
	"github.com/philippebeaulieu/rpi-thermostat/sensor/ds18b20"
	"github.com/philippebeaulieu/rpi-thermostat/sensor/fakesensor"
	"github.com/philippebeaulieu/rpi-thermostat/thermostat"
	"github.com/philippebeaulieu/rpi-thermostat/weather/apixu"
)

func main() {

	var debug = false
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		debug = true
	}

	var err error

	var controller controller.Controller
	if debug {
		controller, err = fakecontroller.NewFakeController()
	} else {
		controller, err = rpi.NewRpiController()
	}

	if err != nil {
		fmt.Printf("failed to create controller: %v\n", err)
		return
	}

	var sensor sensor.Sensor
	if debug {
		sensor, err = fakesensor.NewFakeSensor()
	} else {
		sensor, err = ds18b20.NewDs18b20("28-041685fc45ff")
	}

	if err != nil {
		fmt.Printf("failed to create sensor: %v\n", err)
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

	var server *apiserver.Apiserver
	if debug {
		server = apiserver.NewAPIServer(thermostat, datagatherer, 8080)
	} else {
		server = apiserver.NewAPIServer(thermostat, datagatherer, 80)
	}

	go server.Run()

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
