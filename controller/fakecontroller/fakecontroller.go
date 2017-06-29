package fakecontroller

import "../../controller"

type fakecontroller struct {
}

func NewFakeController() (controller.Controller, error) {
	return &fakecontroller{}, nil
}

func (c *fakecontroller) Off(output int) {
}

func (c *fakecontroller) Heat(output int) {
}
