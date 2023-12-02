package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func DBConfig() *sql.DB {
	// DB connection string
	connectionString := os.Getenv("CONNECTION_STRING")

	if connectionString == "" {
		log.Fatal("Connection string not provided")
	}

	// Initiating DB connection
	db, err := sql.Open("postgres", connectionString)

	if err != nil { // Database connection failed
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil { // Can't ping database
		log.Fatal(err)
	} else { // Database Connected
		fmt.Println("Database Connection Established Successfully")
	}

	return db
}
