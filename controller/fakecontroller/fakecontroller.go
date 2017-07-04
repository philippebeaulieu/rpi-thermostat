package fakecontroller

import "github.com/philippebeaulieu/rpi-thermostat/controller"

type fakecontroller struct {
}

// NewFakeController is use as a constructor
func NewFakeController() (controller.Controller, error) {
	return &fakecontroller{}, nil
}

func (c *fakecontroller) Off(output int) {
}

func (c *fakecontroller) Heat(output int) {
}
