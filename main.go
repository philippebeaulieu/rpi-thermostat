package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stianeikeland/go-rpio"
)

var (
	current = 0
	desired = 21
	sysmode = "off"

	power = 0 //max = 9

	pwm1Total = 0
	pwm2Total = -3
	pwm3Total = -6

	pwm1out = 0
	pwm2out = 0
	pwm3out = 0

	tempSensor = "/sys/bus/w1/devices/28-041685fc45ff/w1_slave"
	logFolder  = "/var/log/rpi-thermostat/"
	logFile    = logFolder + "temp.log"

	heat1 = rpio.Pin(17)
	heat2 = rpio.Pin(21)
	heat3 = rpio.Pin(22)
)

type thermostat struct {
	Current int    `json:"current"`
	Desired int    `json:"desired"`
	Sysmode string `json:"sysmode"`
	Power   int    `json:"power"`
}

func updateHandler(r *http.Request) int {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest
	}

	var t thermostat
	err = json.Unmarshal(body, &t)
	if err != nil {
		return http.StatusBadRequest
	}

	if t.Desired <= 26 {
		desired = t.Desired
	}
	sysmode = t.Sysmode
	return http.StatusOK
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var code int

	switch r.Method {
	case "GET":
		code = http.StatusOK
	case "POST":
		code = updateHandler(r)
	default:
		code = http.StatusNotImplemented
	}
	//log.Printf("current: %v, desired: %v, sysmode: %v, power: %v\n", current, desired, sysmode, power)
	response := thermostat{current, desired, sysmode, power}
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if code != 200 {
		http.Error(w, "", code)
	}
	w.Write(json)
}

// GPIO

func start(p rpio.Pin) {
	p.High()
}

func stop(p rpio.Pin) {
	p.Low()
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

func readTemp() int {
	f, err := os.Open(tempSensor)
	if err != nil {
		log.Println("error opening file= ", err)
		os.Exit(1)
	}

	r := bufio.NewReader(f)
	s, e := readln(r)
	if e != nil {
		log.Println("error opening file= ", err)
		os.Exit(1)
	}
	s, e = readln(r)
	if e != nil {
		log.Println("error opening file= ", err)
		os.Exit(1)
	}

	i, err := strconv.Atoi(s[len(s)-5:])

	temp := i / 100

	return temp
}

func updateOutput(value int, output rpio.Pin) {
	if value == 0 {
		stop(output)
	} else {
		start(output)
	}
}

func pwm(total *int, out *int, output rpio.Pin) {

	*total = *total + 1

	if *total > 8 {
		*total = 0
	}

	if *total >= 0 && *total < power {
		*out = 1
	} else {
		*out = 0
	}

	updateOutput(*out, output)
}

func limitToRange(value, min, max int) int {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

func setupLog() {
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		os.Mkdir(logFolder, 0775)
	}

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)
}

func setupGpio() {
	err := rpio.Open()
	if err != nil {
		log.Printf("Failed to open controller, err: %s\n", err)
		return
	}

	for _, p := range []rpio.Pin{heat1, heat2, heat3} {
		p.Output()
		stop(p)
	}
}

func setupAPIServer() {
	http.HandleFunc("/api", apiHandler)
	http.Handle("/", http.FileServer(http.Dir("/usr/local/rpi-thermostat/src/github.com/philippebeaulieu/rpi-thermostat/ui")))
	go http.ListenAndServe(":80", nil)
}

func updatePower() {
	current = readTemp()

	if sysmode == "off" {
		power = 0
	} else {
		powerSetPoint := limitToRange((desired*10)-current, 0, 9)
		if powerSetPoint > power {
			power = power + 1
		} else if powerSetPoint < power {
			power = power - 1
		}
	}
}

func saveData() {
	db, err := sql.Open("mysql", "thermostat:GDeWFE8Hg3aKh44@tcp(192.168.2.41:3306)/rpi-thermostat?charset=utf8")
	checkErr(err)
	stmt, err := db.Prepare("INSERT temp_data SET time=NOW(),current=?,desired=?,power=?,sysmode=?")
	checkErr(err)
	_, err = stmt.Exec(current, desired, power, sysmode)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func saveDataLoop() {
	for {
		time.Sleep(1 * time.Minute)
		go saveData()
	}
}

func main() {

	setupAPIServer()
	setupLog()
	setupGpio()
	go saveDataLoop()

	for range time.Tick(10 * time.Second) {
		pwm(&pwm1Total, &pwm1out, heat1)
		pwm(&pwm2Total, &pwm2out, heat2)
		pwm(&pwm3Total, &pwm3out, heat3)

		updatePower()

		log.Printf("%v	%v	%v	%v\n", current, desired, power, sysmode)
	}
}
