package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/philippebeaulieu/rpi-thermostat/datagatherer"
	"github.com/philippebeaulieu/rpi-thermostat/thermostat"
)

// Apiserver is use as a reference struct for constructor
type Apiserver struct {
	thermostat   *thermostat.Thermostat
	port         int
	datagatherer *datagatherer.Datagatherer
}

func (s *Apiserver) updateHandler(r *http.Request) int {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusBadRequest
	}
	state := thermostat.State{}
	err = json.Unmarshal(body, &state)
	if err != nil {
		return http.StatusBadRequest
	}
	s.thermostat.Put(state)
	return http.StatusOK
}

func (s *Apiserver) apiHandler(w http.ResponseWriter, r *http.Request) {
	var code int

	switch r.Method {
	case "GET":
		code = http.StatusOK
	case "POST":
		code = s.updateHandler(r)
	default:
		code = http.StatusNotImplemented
	}
	response := s.thermostat.Get()
	// fmt.Printf("%#v\n", response)
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
	// log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}

func (s *Apiserver) apiPastStatesHandler(w http.ResponseWriter, r *http.Request) {
	var code int

	switch r.Method {
	case "GET":
		code = http.StatusOK
	default:
		code = http.StatusNotImplemented
	}
	states := s.datagatherer.GetPreviousDayStates()

	response := convertStatesToReponse(states)

	// fmt.Printf("%#v\n", response)
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
	// log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}

//PastState is used to return past states to webui
type pastState struct {
	Time     []string `json:"time"`
	Interior []int    `json:"interior"`
	Exterior []int    `json:"exterior"`
	SetPoint []int    `json:"setPoint"`
	Power    []int    `json:"power"`
}

func convertStatesToReponse(states []thermostat.State) pastState {
	time := make([]string, 24)
	interior := make([]int, 24)
	exterior := make([]int, 24)
	setPoint := make([]int, 24)
	power := make([]int, 24)

	for i, state := range states {
		time[i] = state.Time.Format("2006-01-02T15:04:05")
		interior[i] = int(state.Current)
		exterior[i] = int(state.Outside.Temp)
		setPoint[i] = state.SetPoint
		power[i] = state.Power
	}

	return pastState{
		Time:     time,
		Interior: interior,
		Exterior: exterior,
		SetPoint: setPoint,
		Power:    power,
	}
}

// NewAPIServer is use as a constructor
func NewAPIServer(thermostat *thermostat.Thermostat, datagatherer *datagatherer.Datagatherer, port int) *Apiserver {
	return &Apiserver{
		thermostat:   thermostat,
		port:         port,
		datagatherer: datagatherer,
	}
}

// Run starts the api processes
func (s *Apiserver) Run() {
	http.HandleFunc("/api", s.apiHandler)
	http.HandleFunc("/api/paststates", s.apiPastStatesHandler)
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}
