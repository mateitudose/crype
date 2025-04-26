package server

import (
	"crype/server/config"
	"database/sql"
	"fmt"
	"log"
	"net"

	pb "crype/api/gen_proto"

	"google.golang.org/grpc"
)

type CrypeServer struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

func NewCrypeServer(db *sql.DB) *CrypeServer {
	return &CrypeServer{db: db}
}

func SetupServer(config *config.ServerConfig) error {
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", config.Port, err)
	}

	s := grpc.NewServer()

	// Register our service implementation
	crypeServer := NewCrypeServer(config.DB)
	pb.RegisterOrderServiceServer(s, crypeServer)

	log.Printf("Server listening on :%s", config.Port)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
