package sensor

type Sensor interface {
	GetTemperature() (float32, error)
}
