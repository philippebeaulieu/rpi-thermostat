package fakesensor

import "github.com/philippebeaulieu/rpi-thermostat/sensor"

type fakesensor struct {
}

// NewFakeSensor is use as a constructor
func NewFakeSensor() (sensor.Sensor, error) {
	return &fakesensor{}, nil
}

func (d *fakesensor) GetTemperature() (float32, error) {
	return 21.2, nil
}
