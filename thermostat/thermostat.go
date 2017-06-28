package thermostat

import (
	"fmt"
	"time"

	"../controller"
	"../sensor"
)

type Thermostat struct {
	sensor     sensor.Sensor
	controller controller.Controller
	sysmode    string
	power      int
	desired    int
	current    int
	pwmTotals  [3]int
}

type State struct {
	Current int    `json:"current"`
	Desired int    `json:"desired"`
	Sysmode string `json:"sysmode"`
	Power   int    `json:"power"`
}

func NewThermostat(sensor sensor.Sensor, controller controller.Controller, desired int) (*Thermostat, error) {
	return &Thermostat{
		sensor:     sensor,
		controller: controller,
		sysmode:    "off",
		desired:    desired,
		current:    0,
		pwmTotals:  [3]int{0, -3, -6},
	}, nil
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
	t.current = state.Current
	t.desired = state.Desired
	t.sysmode = state.Sysmode
	t.Update()
}

func (t *Thermostat) Get() State {
	return State{
		Current: t.current,
		Desired: t.desired,
		Sysmode: t.sysmode,
	}
}

func (t *Thermostat) Update() {
	if t.sysmode == "off" {
		t.power = 0
		t.controller.Off(1)
		t.controller.Off(2)
		t.controller.Off(3)
	} else {

		power := (t.desired * 10) - t.current //current needs to be divided by 10 to get real value IE. 221 = 22.1°C

		if power < 0 {
			power = 0
		}

		if power > 9 {
			power = 9
		}

		t.power = power

	}

	pwm(t, 1)
	pwm(t, 2)
	pwm(t, 3)

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
