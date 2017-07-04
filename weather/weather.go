package weather

// State contains a snapshot of current weather
type State struct {
	TempC     float32
	WindKph   float32
	Humidity  int
	Localtime string
}

// Weather is an interface that allows to get actual weather information in a simplified way
type Weather interface {
	GetWeather() (State, error)
}
