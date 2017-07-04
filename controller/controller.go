package controller

// Controller is an interface that allows to control output in a simplified way
type Controller interface {
	Off(output int)
	Heat(output int)
}
