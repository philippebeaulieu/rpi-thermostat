package thermostat

import (
	"fmt"
	"time"

	"github.com/philippebeaulieu/rpi-thermostat/controller"
	"github.com/philippebeaulieu/rpi-thermostat/sensor"
	"github.com/philippebeaulieu/rpi-thermostat/weather"
)

type Thermostat struct {
	sensor     sensor.Sensor
	controller controller.Controller
	sysmode    string
	power      int
	desired    int
	current    float32
	pwmTotals  [3]int
	Weather    weather.State
}

type State struct {
	Current     float32 `json:"current"`
	Desired     int     `json:"desired"`
	Sysmode     string  `json:"sysmode"`
	Power       int     `json:"power"`
	OutsideTemp float32 `json:"outside_temp"`
	Wind        float32 `json:"wind"`
	Humidity    int     `json:"humidity"`
}

func NewThermostat(sensor sensor.Sensor, controller controller.Controller, desired int) *Thermostat {
	return &Thermostat{
		sensor:     sensor,
		controller: controller,
		sysmode:    "off",
		desired:    desired,
		current:    0,
		pwmTotals:  [3]int{0, -3, -6},
	}
}

func (t *Thermostat) Run() {
	for {
		current, err := t.sensor.GetTemperature()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("switching system off due to sensor failure")
			t.sysmode = "off"
			t.Update()
		} else {
			t.current = current
			t.Update()
		}
		<-time.After(10 * time.Second)
	}
}

func (t *Thermostat) Put(state State) {
	if state.Desired < 5 {
		t.desired = 5
	} else if state.Desired > 30 {
		t.desired = 30
	} else {
		t.desired = state.Desired
	}
	t.sysmode = state.Sysmode
	t.Update()
}

func (t *Thermostat) Get() State {
	return State{
		Current:     t.current,
		Desired:     t.desired,
		Sysmode:     t.sysmode,
		Power:       t.power,
		OutsideTemp: t.Weather.TempC,
		Wind:        t.Weather.WindKph,
		Humidity:    t.Weather.Humidity,
	}
}

func (t *Thermostat) Update() {
	if t.sysmode == "off" {
		t.power = 0
		t.controller.Off(0)
		t.controller.Off(1)
		t.controller.Off(2)
	} else {

		power := int((float32(t.desired) - t.current) * 10)

		if power < 0 {
			power = 0
		}

		if power > 9 {
			power = 9
		}

		t.power = power

	}

	pwm(t, 0)
	pwm(t, 1)
	pwm(t, 2)
}

func pwm(t *Thermostat, output int) {

	t.pwmTotals[output] = t.pwmTotals[output] + 1

	if t.pwmTotals[output] > 8 {
		t.pwmTotals[output] = 0
	}

	if t.pwmTotals[output] >= 0 && t.pwmTotals[output] < t.power {
		t.controller.Heat(output)
	} else {
		t.controller.Off(output)
	}

}
