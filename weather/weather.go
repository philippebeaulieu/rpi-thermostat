package weather

type State struct {
	TempC     float32
	WindKph   float32
	Humidity  int
	Localtime string
}

type Weather interface {
	GetWeather() (State, error)
}
