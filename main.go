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
}
