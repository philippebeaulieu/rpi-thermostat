package thermostat

import (
	"fmt"
	"time"

	"../controller"
	"../sensor"
)

var (
	power = 0

	pwm1Total = 0
	pwm2Total = -3
	pwm3Total = -6

	pwm1out = 0
	pwm2out = 0
	pwm3out = 0
)

type Thermostat struct {
	sensor     sensor.Sensor
	controller controller.Controller
	sysmode    string
	power      int
	desired    int
	current    int
}

type ThermostatState struct {
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
	}, nil
}

func (t *Thermostat) Run() {
	for {
		current, err := t.sensor.GetTemperature()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("switching system off due to sensor failure")
			t.controller.Off(1)
			t.controller.Off(2)
			t.controller.Off(3)
			t.sysmode = "off"
		} else {
			t.current = current
			t.Update()
		}
		<-time.After(10 * time.Second)
	}
}

func (t *Thermostat) Put(state ThermostatState) {
	t.current = state.Current
	t.desired = state.Desired
	t.sysmode = state.Sysmode
	t.Update()
}

func (t *Thermostat) Get() ThermostatState {
	return ThermostatState{
		Current: t.current,
		Desired: t.desired,
		Sysmode: t.sysmode,
	}
}

func (t *Thermostat) Update() {
	if t.sysmode == "off" {
		power = 0
		t.controller.Off(1)
		t.controller.Off(2)
		t.controller.Off(3)
	} else {
		adjustPower(t.desired, t.current)
	}

	pwm(t, &pwm1Total, &pwm1out, 1)
	pwm(t, &pwm2Total, &pwm2out, 2)
	pwm(t, &pwm3Total, &pwm3out, 3)

}

func adjustPower(desired int, current int) {
	powerSetPoint := limitToRange((desired*10)-current, 0, 9)
	if powerSetPoint > power {
		power = power + 1
	} else if powerSetPoint < power {
		power = power - 1
	}
}

func limitToRange(value, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

func pwm(t *Thermostat, total *int, out *int, output int) {

	*total = *total + 1

	if *total > 8 {
		*total = 0
	}

	if *total >= 0 && *total < power {
		*out = 1
	} else {
		*out = 0
	}

	if *out == 0 {
		t.controller.Off(output)
	} else {
		t.controller.Heat(output)
	}
}
