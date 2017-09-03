package fakecontroller

import (
	"fmt"

	"github.com/philippebeaulieu/rpi-thermostat/controller"
)

type fakecontroller struct {
}

// NewFakeController is use as a constructor
func NewFakeController() (controller.Controller, error) {
	return &fakecontroller{}, nil
}

func (c *fakecontroller) Off(output int) {
	fmt.Printf("Pin %v OFF\n", output)
}

func (c *fakecontroller) Heat(output int) {
	fmt.Printf("Pin %v ON\n", output)
}
