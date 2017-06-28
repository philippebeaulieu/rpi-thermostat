package rpi

import (
	"fmt"

	"../../controller"
	"github.com/stianeikeland/go-rpio"
)

type rpiController struct {
	heat1 rpio.Pin
	heat2 rpio.Pin
	heat3 rpio.Pin
}

func NewRpiController() (controller.Controller, error) {
	err := rpio.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open RPi controller: %s", err)
	}
	return &rpiController{
		heat1: rpio.Pin(17),
		heat2: rpio.Pin(21),
		heat3: rpio.Pin(22),
	}, nil
}

func (c *rpiController) Off(output int) {
	if output == 1 {
		c.heat1.High()
	} else if output == 2 {
		c.heat2.High()
	} else if output == 3 {
		c.heat3.High()
	}
}

func (c *rpiController) Heat(output int) {
	if output == 1 {
		c.heat1.Low()
	} else if output == 2 {
		c.heat2.Low()
	} else if output == 3 {
		c.heat3.Low()
	}
}
