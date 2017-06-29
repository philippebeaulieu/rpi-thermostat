package fake

import "../../controller"

type fakeController struct {
}

func NewFakeController() (controller.Controller, error) {
	return &fakeController{}, nil
}

func (c *fakeController) Off(output int) {
}

func (c *fakeController) Heat(output int) {
}
