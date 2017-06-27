package controller

type Controller interface {
	Off(output int)
	Heat(output int)
}
