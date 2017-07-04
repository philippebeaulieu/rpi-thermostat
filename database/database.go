package database

import (
	"database/sql"

	"github.com/philippebeaulieu/rpi-thermostat/thermostat"

	// forced import to provide actual mysql implementation without having to directly refere to it in code
	_ "github.com/go-sql-driver/mysql"
	"github.com/philippebeaulieu/rpi-thermostat/thermostat/statequeue"
)

// Database is use as a reference struct for constructor
type Database struct {
	thermostat *thermostat.Thermostat
	db         *sql.DB
}

// NewDatabase is use as a constructor
func NewDatabase(thermostat *thermostat.Thermostat) (*Database, error) {
	db, err := sql.Open("mysql", "thermostat:GDeWFE8Hg3aKh44@tcp(192.168.2.41:3306)/rpi-thermostat?charset=utf8")
	if err != nil {
		return nil, err
	}

	return &Database{
		thermostat: thermostat,
		db:         db,
	}, nil
}

// GetPastStates returns a list of the last days saved states
func (d *Database) GetPastStates() ([]thermostat.State, error) {
	rows, err := d.db.Query("SELECT current, desired, sysmode, outside_temp, wind, humidity FROM temp_data WHERE time BETWEEN DATE_SUB(NOW(), INTERVAL 1 DAY) AND NOW() LIMIT 1440; ")
	if err != nil {
		return nil, err
	}

	queue := statequeue.NewQueue(1440)

	for rows.Next() {
		var current int
		var desired int
		var sysmode string
		var outsideTemp int
		var wind int
		var humidity int

		err = rows.Scan(&current, &desired, &sysmode, &outsideTemp, &wind, &humidity)
		if err != nil {
			return nil, err
		}

		state := thermostat.State{
			Current:     float32(current),
			Desired:     desired,
			Sysmode:     sysmode,
			OutsideTemp: float32(outsideTemp),
			Wind:        float32(wind),
			Humidity:    humidity,
		}

		queue.Push(state)
	}

	return queue.ToArray(), nil
}

// SaveData saves data to database
func (d *Database) SaveData(state thermostat.State) error {
	stmt, err := d.db.Prepare("INSERT temp_data SET time=NOW(),current=?,desired=?,power=?,sysmode=?, outside_temp=?, wind=?, humidity=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(state.Current, state.Desired, state.Power, state.Sysmode, int(state.OutsideTemp), int(state.Wind), state.Humidity)

	return err
}
