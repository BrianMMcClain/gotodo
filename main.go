package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {

	// Read in the config
	configPath := flag.String("config", "./config.json", "Path to the JSON configuration file")
	flag.Parse()
	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)
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
	r.DELETE("/:id", deleteToDo(db))
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
