package db

import (
	"database/sql"
	"log"
)

func InitDb() {
	db, err := sql.Open("sqlite3", "ami.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Print("creating todos table if it doesnt exist..")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    msg TEXT,
		done INTEGER
)`)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("creating reminders table if it doesnt exist..")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reminders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT,
		cadence TEXT
)`)
	if err != nil {
		log.Fatal(err)
	}
}
