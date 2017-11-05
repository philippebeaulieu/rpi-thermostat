package fakesensor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/philippebeaulieu/rpi-thermostat/sensor"
)

type fakesensor struct {
}

// NewFakeSensor is use as a constructor
func NewFakeSensor() (sensor.Sensor, error) {
	return &fakesensor{}, nil
}

func (d *fakesensor) GetTemperature() (float32, error) {
	f, err := os.Open("testSensor")
	if err != nil {
		return -1, err
	}

	r := bufio.NewReader(f)
	s, err := readln(r)
	if err != nil {
		return -1, err
	}
	s, err = readln(r)
	if err != nil {
		return -1, err
	}

	fmt.Printf("s: %v\n", s)
	fmt.Printf("subS: %v\n", s[strings.Index(s, "=")+1:])

	temp, err := strconv.Atoi(s[strings.Index(s, "=")+1:])

	return float32(temp) / 1000, err

}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
