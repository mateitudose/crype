package main

import (
	"crype/server"
	"crype/server/config"
	"crype/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Initialize database connection
	db := utils.ConnectDB()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	serverConfig := config.NewServerConfig(os.Getenv("CRYPE_PORT"), db)
	if err := server.SetupServer(serverConfig); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
