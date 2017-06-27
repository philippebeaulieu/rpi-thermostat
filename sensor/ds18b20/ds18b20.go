package ds18b20

import (
	"bufio"
	"os"
	"strconv"

	"../../sensor"
)

var tempSensor = "/sys/bus/w1/devices/28-041685fc45ff/w1_slave"

type ds18b20s struct {
	sensor.Sensor
}

func NewDs18b20(deviceid string) sensor.Sensor {
	s := &ds18b20s{}
	return s
}

func GetTemperature() (int, error) {

	f, err := os.Open(tempSensor)
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
