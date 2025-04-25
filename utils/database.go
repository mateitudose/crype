package utils

import (
	"crype/api/gen_jet/crype_db/public/model"
	pb "crype/api/gen_proto"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// ConnectDB establishes a connection to the PostgreSQL database
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

	log.Println("Connected to PostgreSQL DB!")
	return db
}

// GetProtobufOrderStatus converts a database order status to its protobuf equivalent
func GetProtobufOrderStatus(dbStatus model.OrderStatus) (pb.OrderStatus, bool) {
	orderStatusMap := map[model.OrderStatus]pb.OrderStatus{
		model.OrderStatus_Pending:    pb.OrderStatus_PENDING,
		model.OrderStatus_Processing: pb.OrderStatus_PROCESSING,
		model.OrderStatus_Confirmed:  pb.OrderStatus_CONFIRMED,
		model.OrderStatus_Failed:     pb.OrderStatus_FAILED,
		model.OrderStatus_Canceled:   pb.OrderStatus_CANCELED,
	}

	status, exists := orderStatusMap[dbStatus]
	return status, exists
}
