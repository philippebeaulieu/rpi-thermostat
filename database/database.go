package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func SaveData(current int, desired int, power int, sysmode string) {
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
