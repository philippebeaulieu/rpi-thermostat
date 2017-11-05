package ds18b20

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/philippebeaulieu/rpi-thermostat/sensor"
)

type ds18b20 struct {
	deviceid string
}

// NewDs18b20 is use as a constructor
func NewDs18b20(deviceid string) (sensor.Sensor, error) {

	ds18b20 := &ds18b20{
		deviceid: deviceid,
	}

	_, err := ds18b20.GetTemperature()
	if err != nil {
		return nil, err
	}

	return ds18b20, nil
}

func (d *ds18b20) GetTemperature() (float32, error) {

	f, err := os.Open("/sys/bus/w1/devices/" + d.deviceid + "/w1_slave")
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
