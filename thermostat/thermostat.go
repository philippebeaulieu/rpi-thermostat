package thermostat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/philippebeaulieu/rpi-thermostat/controller"
	"github.com/philippebeaulieu/rpi-thermostat/sensor"
	"github.com/philippebeaulieu/rpi-thermostat/weather"
)

// Thermostat is use as a reference struct for constructor
type Thermostat struct {
	sensor     sensor.Sensor
	controller controller.Controller
	pwmTotals  [3]int
	state      State
}

// State contains a snapshot of current data
type State struct {
	Time     time.Time `json:"time"`
	Current  float32   `json:"current"`
	SetPoint int       `json:"set_point"`
	Sysmode  string    `json:"sysmode"`
	Power    int       `json:"power"`
	Outside  Outside   `json:"outside"`
	Settings Settings  `json:"settings"`
}

// Outside information
type Outside struct {
	Temp     float32 `json:"temp"`
	Wind     float32 `json:"wind"`
	Humidity int     `json:"humidity"`
}

// Settings for thermostat
type Settings struct {
	ManualTemp         int `json:"manual_temp"`
	PresentTemp        int `json:"present_temp"`
	AbsentTemp         int `json:"absent_temp"`
	MondayStartHour    int `json:"monday_start_hour"`
	MondayStopHour     int `json:"monday_stop_hour"`
	TuesdayStartHour   int `json:"tuesday_start_hour"`
	TuesdayStopHour    int `json:"tuesday_stop_hour"`
	WednesdayStartHour int `json:"wednesday_start_hour"`
	WednesdayStopHour  int `json:"wednesday_stop_hour"`
	ThursdayStartHour  int `json:"thursday_start_hour"`
	ThursdayStopHour   int `json:"thursday_stop_hour"`
	FridayStartHour    int `json:"friday_start_hour"`
	FridayStopHour     int `json:"friday_stop_hour"`
	SaturdayStartHour  int `json:"saturday_start_hour"`
	SaturdayStopHour   int `json:"saturday_stop_hour"`
	SundayStartHour    int `json:"sunday_start_hour"`
	SundayStopHour     int `json:"sunday_stop_hour"`
}

// NewThermostat is use as a constructor
func NewThermostat(sensor sensor.Sensor, controller controller.Controller) *Thermostat {
	return &Thermostat{
		sensor:     sensor,
		controller: controller,
		state: State{
			Sysmode:  "off",
			Outside:  Outside{},
			Settings: loadSettings(),
		},
		pwmTotals: [3]int{0, -3, -6},
	}
}

// Run starts the thermostat processes
func (t *Thermostat) Run() {
	for {
		t.state.Time = time.Now()
		current, err := t.sensor.GetTemperature()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("switching system off due to sensor failure")
			t.state.Sysmode = "off"
			t.update()
		} else {
			t.state.Current = current
			t.update()
		}
		<-time.After(10 * time.Second)
	}
}

// Put receives a state an transfer its values to the thermostat
func (t *Thermostat) Put(state State) {
	t.state.Settings.ManualTemp = minMax(state.Settings.ManualTemp, 5, 30)

	t.state.Settings.PresentTemp = minMax(state.Settings.PresentTemp, 5, 30)
	t.state.Settings.AbsentTemp = minMax(state.Settings.AbsentTemp, 5, 30)

	t.state.Settings.MondayStartHour = minMax(state.Settings.MondayStartHour, 0, 23)
	t.state.Settings.MondayStopHour = minMax(state.Settings.MondayStopHour, 0, 23)

	t.state.Settings.TuesdayStartHour = minMax(state.Settings.TuesdayStartHour, 0, 23)
	t.state.Settings.TuesdayStopHour = minMax(state.Settings.TuesdayStopHour, 0, 23)

	t.state.Settings.WednesdayStartHour = minMax(state.Settings.WednesdayStartHour, 0, 23)
	t.state.Settings.WednesdayStopHour = minMax(state.Settings.WednesdayStopHour, 0, 23)

	t.state.Settings.ThursdayStartHour = minMax(state.Settings.ThursdayStartHour, 0, 23)
	t.state.Settings.ThursdayStopHour = minMax(state.Settings.ThursdayStopHour, 0, 23)

	t.state.Settings.FridayStartHour = minMax(state.Settings.FridayStartHour, 0, 23)
	t.state.Settings.FridayStopHour = minMax(state.Settings.FridayStopHour, 0, 23)

	t.state.Settings.SaturdayStartHour = minMax(state.Settings.SaturdayStartHour, 0, 23)
	t.state.Settings.SaturdayStopHour = minMax(state.Settings.SaturdayStopHour, 0, 23)

	t.state.Settings.SundayStartHour = minMax(state.Settings.SundayStartHour, 0, 23)
	t.state.Settings.SundayStopHour = minMax(state.Settings.SundayStopHour, 0, 23)

	saveSettings(t.state.Settings)

	t.state.Sysmode = state.Sysmode

	t.update()
}

func minMax(val int, min int, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

// Get takes actual thermostat values and returns them as a State
func (t *Thermostat) Get() State {
	return t.state
}

func (t *Thermostat) update() {

	if t.state.Sysmode == "automatic" {
		t.state.SetPoint = getAutomaticTemp(t.state)
		applyPower(t)
	} else if t.state.Sysmode == "manual" {
		t.state.SetPoint = t.state.Settings.ManualTemp
		applyPower(t)
	} else if t.state.Sysmode == "off" {
		t.state.Power = 0
		t.controller.Off(0)
		t.controller.Off(1)
		t.controller.Off(2)
	}
}

func applyPower(t *Thermostat) {

	var MinPower = 0
	var MaxPower = 9

	var setPoint = float64(t.state.SetPoint)
	var current = float64(t.state.Current)
	var outside = float64(t.state.Outside.Temp)

	power := int(math.Ceil((setPoint - current) * ((outside - setPoint - 10) / -10.0)))

	if power < MinPower {
		power = MinPower
	}

	if power > MaxPower {
		power = MaxPower
	}

	t.state.Power = power

	pwm(t, 0)
	pwm(t, 1)
	pwm(t, 2)
}

func getAutomaticTemp(state State) int {
	switch state.Time.Weekday() {
	case 0:
		return getTempForRange(state.Time.Hour(), state.Settings.SundayStartHour, state.Settings.SundayStopHour, state.Settings.PresentTemp, state.Settings.AbsentTemp)
	case 1:
		return getTempForRange(state.Time.Hour(), state.Settings.MondayStartHour, state.Settings.MondayStopHour, state.Settings.PresentTemp, state.Settings.AbsentTemp)
	case 2:
		return getTempForRange(state.Time.Hour(), state.Settings.TuesdayStartHour, state.Settings.TuesdayStopHour, state.Settings.PresentTemp, state.Settings.AbsentTemp)
	case 3:
		return getTempForRange(state.Time.Hour(), state.Settings.WednesdayStartHour, state.Settings.WednesdayStopHour, state.Settings.PresentTemp, state.Settings.AbsentTemp)
	case 4:
		return getTempForRange(state.Time.Hour(), state.Settings.ThursdayStartHour, state.Settings.ThursdayStopHour, state.Settings.PresentTemp, state.Settings.AbsentTemp)
	case 5:
		return getTempForRange(state.Time.Hour(), state.Settings.FridayStartHour, state.Settings.FridayStopHour, state.Settings.PresentTemp, state.Settings.AbsentTemp)
	case 6:
		return getTempForRange(state.Time.Hour(), state.Settings.SaturdayStartHour, state.Settings.SaturdayStopHour, state.Settings.PresentTemp, state.Settings.AbsentTemp)
	}

	return -1
}

func getTempForRange(hour int, startHour int, stopHour int, presentTemp int, absentTemp int) int {
	if hour >= startHour && hour < stopHour {
		return presentTemp
	}

	return absentTemp
}

func pwm(t *Thermostat, output int) {

	t.pwmTotals[output] = t.pwmTotals[output] + 1

	if t.pwmTotals[output] > 8 {
		t.pwmTotals[output] = 0
	}

	if t.pwmTotals[output] >= 0 && t.pwmTotals[output] < t.state.Power {
		t.controller.Heat(output)
	} else {
		t.controller.Off(output)
	}

}

// LoadWeatherState receives a weather state and loads it into the thermostat state
func (t *Thermostat) LoadWeatherState(weather weather.State) {
	t.state.Outside.Temp = weather.TempC
	t.state.Outside.Humidity = weather.Humidity
	t.state.Outside.Wind = weather.WindKph
}

func loadSettings() Settings {
	val, err := ioutil.ReadFile("settings.json")
	if err != nil {
		defaultSettings := defaultSettings()
		saveSettings(defaultSettings)
		return defaultSettings
	}

	settings := Settings{}
	err = json.Unmarshal(val, &settings)
	if err != nil {
		defaultSettings := defaultSettings()
		saveSettings(defaultSettings)
		return defaultSettings
	}

	return settings
}

func defaultSettings() Settings {
	return Settings{
		ManualTemp:         21,
		PresentTemp:        21,
		AbsentTemp:         10,
		MondayStartHour:    0,
		MondayStopHour:     0,
		TuesdayStartHour:   0,
		TuesdayStopHour:    0,
		WednesdayStartHour: 0,
		WednesdayStopHour:  0,
		ThursdayStartHour:  0,
		ThursdayStopHour:   0,
		FridayStartHour:    16,
		FridayStopHour:     23,
		SaturdayStartHour:  7,
		SaturdayStopHour:   23,
		SundayStartHour:    7,
		SundayStopHour:     23,
	}
}

func saveSettings(settings Settings) {
	json, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("settings.json", json, 0644)
	if err != nil {
		panic(err)
	}
}
