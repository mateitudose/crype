package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func ConnectDB() *sql.DB {
	// Load environment variables from .env file
	user := os.Getenv("CRYPE_DB_USER")
	password := os.Getenv("CRYPE_DB_PASSWORD")
	dbName := os.Getenv("CRYPE_DB_NAME")
	host := os.Getenv("CRYPE_DB_HOST")
	port := os.Getenv("CRYPE_DB_PORT")

	// TODO: In the future, an external PostgreSQL database could be used
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// If there is an error connecting to the database, log the fatal error and exit the program
		log.Fatal(err)
	}

	// Test the connection by pinging the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Connected to PostgreSQL DB!")
	return db
}
