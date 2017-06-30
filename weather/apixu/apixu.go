package apixu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/philippebeaulieu/rpi-thermostat/thermostat"
	"github.com/philippebeaulieu/rpi-thermostat/weather"
)

type Apixu struct {
	thermostat *thermostat.Thermostat
	apikey     string
	location   string
}

type apixuLocationReponse struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float32 `json:"lat"`
	Lon            float32 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}
type apixuCurrentReponse struct {
	TempC    float32 `json:"temp_c"`
	WindKph  float32 `json:"wind_kph"`
	Humidity int     `json:"humidity"`
}

type apixuResponse struct {
	Location apixuLocationReponse `json:"location"`
	Current  apixuCurrentReponse  `json:"current"`
}

func NewApixuWeather(thermostat *thermostat.Thermostat, apikey string, location string) *Apixu {
	return &Apixu{
		thermostat: thermostat,
		apikey:     apikey,
		location:   location,
	}
}

func (a *Apixu) GetWeather() (weather.State, error) {
	res, err := http.Get("https://api.apixu.com/v1/current.json?key=" + a.apikey + "&q=" + a.location)
	if err != nil {
		return weather.State{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return weather.State{}, err
	}

	var s = new(apixuResponse)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return weather.State{}, err
	}

	return weather.State{
		Localtime: s.Location.Localtime,
		TempC:     s.Current.TempC,
		Humidity:  s.Current.Humidity,
		WindKph:   s.Current.WindKph,
	}, nil
}

func (a *Apixu) Run() {
	for {
		weather, err := a.GetWeather()
		if err == nil {
			a.thermostat.Weather = weather
		}
		<-time.After(15 * time.Minute)
	}
}
