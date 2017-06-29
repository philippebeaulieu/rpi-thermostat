package fakesensor

import "../../sensor"

type fakesensor struct {
}

func NewFakeSensor() (sensor.Sensor, error) {
	return &fakesensor{}, nil
}

func (d *fakesensor) GetTemperature() (int, error) {
	return 21, nil
}
