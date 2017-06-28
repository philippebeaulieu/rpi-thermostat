package ds18b20

import (
	"bufio"
	"os"
	"strconv"

	"../../sensor"
)

type ds18b20 struct {
	deviceid string
}

func NewDs18b20(deviceid string) (sensor.Sensor, error) {
	return &ds18b20{
		deviceid: deviceid,
	}, nil
}

func (d *ds18b20) GetTemperature() (int, error) {

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

	temp, err := strconv.Atoi(s[len(s)-5:])

	return temp / 100, err
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
