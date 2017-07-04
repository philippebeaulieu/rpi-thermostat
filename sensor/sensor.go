package sensor

// Sensor is an interface that allows to get actual temperature information in a simplified way
type Sensor interface {
	GetTemperature() (float32, error)
}
