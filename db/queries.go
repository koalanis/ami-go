package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func TestDb() {
	log.Print("Hlelo")
	log.Print(os.Getwd())
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

	_, err = db.Exec(`INSERT OR IGNORE INTO reminders (id, topic, cadence) VALUES (1, ?, ?)`, "leetcode", "DAILY")
	if err != nil {
		log.Fatal(err)
	}

	// Query data
	rows, err := db.Query("SELECT id, msg FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var msg string
		if err := rows.Scan(&id, &msg); err != nil {
			log.Fatal(err)
		}
		log.Printf("ID: %d, Msg: %s,", id, msg)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("SELECT id, topic, cadence FROM reminders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var topic, cadence string
		if err := rows.Scan(&id, &topic, &cadence); err != nil {
			log.Fatal(err)
		}
		log.Printf("ID: %d, topic: %s, cadence: %s", id, topic, cadence)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func GetTodos() string {
	log.Print("Hlelo")
	log.Print(os.Getwd())
	db, err := sql.Open("sqlite3", "ami.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query data
	rows, err := db.Query("SELECT id, msg FROM todos WHERE done = 0")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	collectedTodos := make([]string, 0)

	for rows.Next() {
		var id int
		var msg string
		if err := rows.Scan(&id, &msg); err != nil {
			log.Fatal(err)
		}
		collectedTodos = append(collectedTodos, fmt.Sprintf("ID: %d Task: %s", id, msg))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	if len(collectedTodos) == 0 {
		return "[]"
	}
	return strings.Join(collectedTodos, "\n")
}

func AddTodo(msg string) {
	db, err := sql.Open("sqlite3", "ami.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert data
	_, err = db.Exec(`INSERT INTO todos (msg, done) VALUES (?, ?)`, msg, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func DoTodo(id string) {
	db, err := sql.Open("sqlite3", "ami.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Insert data
	_, err = db.Exec(`UPDATE OR IGNORE todos SET done = 1 WHERE id = ?`, intId, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func ClearTodos() {
	db, err := sql.Open("sqlite3", "ami.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Insert data
	_, err = db.Exec(`DELETE FROM todos`)
	if err != nil {
		log.Fatal(err)
	}
}
