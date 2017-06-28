package database

import (
	"database/sql"
	"time"

	"../thermostat"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	thermostat *thermostat.Thermostat
}

func NewDatabase(thermostat *thermostat.Thermostat) (*Database, error) {
	return &Database{
		thermostat: thermostat,
	}, nil
}

func (d *Database) Run() {
	for {
		saveData(d.thermostat.Get().Current, d.thermostat.Get().Desired, d.thermostat.Get().Power, d.thermostat.Get().Sysmode)
		<-time.After(1 * time.Minute)
	}
}

func saveData(current int, desired int, power int, sysmode string) {
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
