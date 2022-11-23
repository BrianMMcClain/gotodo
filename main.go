package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "pass"
	dbname   = "gotodo"
)

type ToDo struct {
	Id   int
	Text string
	Done bool
}

func main() {

	// Connect to the database
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

	// If needed, reinit the database
	if os.Getenv("INITDB") == "true" {
		initDB(db)
	}

	// Configure Gin
	r := gin.Default()
	r.GET("/", GetToDos(db))
	r.Run(":8080")
}

func GetToDos(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM todos")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var todos []ToDo
		for rows.Next() {
			var todo ToDo
			rows.Scan(&todo.Id, &todo.Text, &todo.Done)
			todos = append(todos, todo)
		}

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, todos)
	}

	return gin.HandlerFunc(fn)
}

func initDB(db *sql.DB) {
	log.Printf("Initializing database . . .")

	// Recreate the table
	db.Exec("DROP TABLE IF EXISTS todos;")
	createTableQuery := "CREATE TABLE IF NOT EXISTS todos (id serial, text text, done BOOLEAN DEFAULT false, PRIMARY KEY (id));"
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Dummy data
	insertStr := "INSERT INTO todos(id, text, done) VALUES ($1, $2, $3);"
	db.Exec(insertStr, 1, "Take out the trash", false)
	db.Exec(insertStr, 2, "Do the dishes", false)
	db.Exec(insertStr, 3, "Mop the floors", true)
}
