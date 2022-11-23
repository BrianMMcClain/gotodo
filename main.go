package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "pass"
	dbname   = "gotodo"
)

var db *sql.DB

var ToDo struct {
	Id   int
	Text string
	Done bool
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Close the connection to the DB when the server closes
	defer db.Close()

	// Ensure the connection to the database was created properly
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	initDB()
}

func initDB() {
	createStr := `CREATE TABLE IF NOT EXISTS todos(
		id INT PRIMARY  KEY         NOT NULL,
		text            TEXT        NOT NULL,
		done            BOOLEAN     NOT NULL,
	 );`

	fmt.Println(createStr)

	_, err := db.Exec(createStr)
	if err != nil {
		log.Fatal(err)
	}

	insertStr := "INSERT INTO todos(id, text, done) VALUES (?, ?, ?)"
	db.Exec(insertStr, 1, "Take out the trash", false)
	db.Exec(insertStr, 2, "Do the dishes", false)
	db.Exec(insertStr, 3, "Mop the floors", true)

	row := db.QueryRow("SELECT * FROM todos WHERE id=3 LIMIT 1")
	fmt.Println(row)
}
