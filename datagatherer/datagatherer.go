package datagatherer

import (
	"fmt"
	"time"

	"github.com/philippebeaulieu/rpi-thermostat/database"
	"github.com/philippebeaulieu/rpi-thermostat/thermostat"
	"github.com/philippebeaulieu/rpi-thermostat/thermostat/statequeue"
)

// Datagatherer is use as a reference struct for constructor
type Datagatherer struct {
	database   *database.Database
	thermostat *thermostat.Thermostat
	queue      *statequeue.Queue
}

// PastDayCost contains crunched data from past days
type PastDayCost struct {
	Datetime string `json:"datetime"`
	KiloWatt int    `json:"kilowatt"`
}

// NewDataGatherer is use as a constructor
func NewDataGatherer(thermostat *thermostat.Thermostat, database *database.Database) (*Datagatherer, error) {
	q := statequeue.NewQueue(1440)
	states, err := database.GetPastStates()
	if err != nil {
		return nil, err
	}

	for _, state := range states {
		q.Push(state)
	}

	return &Datagatherer{
		database:   database,
		thermostat: thermostat,
		queue:      q,
	}, nil
}

// GetPreviousDayStates returns an array of the previous day states (one for every hour)
func (d *Datagatherer) GetPreviousDayStates() []thermostat.State {
	states := d.queue.ToArray()
	exportedStates := make([]thermostat.State, 24)
	exportedStatesIndex := 0

	for index, state := range states {
		if index%60 == 0 {
			exportedStates[exportedStatesIndex] = state
			exportedStatesIndex++
		}
	}

	return exportedStates
}

// func (d *Datagatherer) getPreviousWeekCosts() ([]PastDayCost, error) {
// 	return nil, nil
// }

// Run starts the datagatherer processes
func (d *Datagatherer) Run() {
	for {
		state := d.thermostat.Get()
		state.Time = time.Now()
		err := d.database.SaveData(state)
		if err != nil {
			fmt.Printf("failed to save data: %v\n", err)
		}

		d.queue.Push(state)
		<-time.After(1 * time.Minute)
	}
}
