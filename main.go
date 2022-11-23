package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Id   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func getToDos(db *sql.DB) gin.HandlerFunc {
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

func getToDo(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var todo ToDo
		db.QueryRow("SELECT * FROM todos WHERE id=$1", c.Param("id")).Scan(&todo.Id, &todo.Text, &todo.Done)

		c.Header("Access-Control-Allow-Origin", "*")

		if todo.Id > 0 {
			c.JSON(http.StatusOK, todo)
		} else {
			c.Status(http.StatusNotFound)
		}
	}

	return gin.HandlerFunc(fn)
}

func addToDo(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatal(err)
		}

		var todo ToDo
		json.Unmarshal(jsonData, &todo)

		queryString := "INSERT INTO todos(text, done) VALUES ($1, $2) RETURNING id"
		var id int
		db.QueryRow(queryString, todo.Text, todo.Done).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}

		c.String(http.StatusOK, fmt.Sprint(id))
	}

	return gin.HandlerFunc(fn)
}

func updateToDo(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatal(err)
		}

		var todo ToDo
		json.Unmarshal(jsonData, &todo)

		var insertError error
		if len(todo.Text) > 0 {
			queryString := "UPDATE todos SET text=$1, done=$2 WHERE id=$3"
			_, insertError = db.Exec(queryString, todo.Text, todo.Done, c.Param("id"))
		} else {
			queryString := "UPDATE todos SET done=$1 WHERE id=$2"
			_, insertError = db.Exec(queryString, todo.Done, c.Param("id"))
		}

		if insertError != nil {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusOK)
		}
	}

	return gin.HandlerFunc(fn)
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
	r.GET("/", getToDos(db))
	r.GET("/:id", getToDo(db))
	r.POST("/", addToDo(db))
	r.POST("/:id", updateToDo(db))
	r.Run(":8080")
}

func initDB(db *sql.DB) {
	log.Printf("Initializing database . . .")

	// Recreate the table
	db.Exec("DROP TABLE IF EXISTS todos;")
	createTableQuery := "CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, text TEXT, done BOOLEAN DEFAULT false);"
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// Dummy data
	insertStr := "INSERT INTO todos(text, done) VALUES ($1, $2);"
	db.Exec(insertStr, "Take out the trash", false)
	db.Exec(insertStr, "Do the dishes", false)
	db.Exec(insertStr, "Mop the floors", true)
}
