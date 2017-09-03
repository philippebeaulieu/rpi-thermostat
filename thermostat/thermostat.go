package thermostat

import (
	"fmt"
	"time"

	"github.com/philippebeaulieu/rpi-thermostat/controller"
	"github.com/philippebeaulieu/rpi-thermostat/sensor"
	"github.com/philippebeaulieu/rpi-thermostat/weather"
)

// Thermostat is use as a reference struct for constructor
type Thermostat struct {
	sensor     sensor.Sensor
	controller controller.Controller
	pwmTotals  [3]int
	state      State
}

// State contains a snapshot of current data
type State struct {
	Time        time.Time `json:"time"`
	Current     float32   `json:"current"`
	Desired     int       `json:"desired"`
	Sysmode     string    `json:"sysmode"`
	Power       int       `json:"power"`
	OutsideTemp float32   `json:"outside_temp"`
	Wind        float32   `json:"wind"`
	Humidity    int       `json:"humidity"`
}

// NewThermostat is use as a constructor
func NewThermostat(sensor sensor.Sensor, controller controller.Controller, desired int) *Thermostat {
	return &Thermostat{
		sensor:     sensor,
		controller: controller,
		state: State{
			Sysmode: "off",
			Desired: desired,
			Current: 0,
		},
		pwmTotals: [3]int{0, -3, -6},
	}
}

// Run starts the thermostat processes
func (t *Thermostat) Run() {
	for {
		current, err := t.sensor.GetTemperature()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("switching system off due to sensor failure")
			t.state.Sysmode = "off"
			t.update()
		} else {
			t.state.Current = current
			t.update()
		}
		<-time.After(10 * time.Second)
	}
}

// Put receives a state an transfer its values to the thermostat
func (t *Thermostat) Put(state State) {
	if state.Desired < 5 {
		t.state.Desired = 5
	} else if state.Desired > 30 {
		t.state.Desired = 30
	} else {
		t.state.Desired = state.Desired
	}
	t.state.Sysmode = state.Sysmode
	t.update()
}

// Get takes actual thermostat values and returns them as a State
func (t *Thermostat) Get() State {
	return t.state
}

func (t *Thermostat) update() {
	if t.state.Sysmode == "off" {
		t.state.Power = 0
		t.controller.Off(0)
		t.controller.Off(1)
		t.controller.Off(2)
	} else {

		power := int((float32(t.state.Desired) - t.state.Current) * 10)

		if power < 0 {
			power = 0
		}

		if power > 9 {
			power = 9
		}

		gain := float32(float32(t.state.Desired)-t.state.OutsideTemp) / 30.0

		if gain < 0.0 {
			gain = 0.0
		}

		if gain > 1.0 {
			gain = 1.0
		}

		t.state.Power = int(float32(power) * gain)

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

	if t.pwmTotals[output] >= 0 && t.pwmTotals[output] < t.state.Power {
		t.controller.Heat(output)
	} else {
		t.controller.Off(output)
	}

}

// LoadWeatherState receives a weather state and loads it into the thermostat state
func (t *Thermostat) LoadWeatherState(weather weather.State) {
	t.state.OutsideTemp = weather.TempC
	t.state.Humidity = weather.Humidity
	t.state.Wind = weather.WindKph
}
