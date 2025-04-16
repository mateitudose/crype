package main

import (
	"context"
	pb "crype/api/gen_proto"
	"crype/utils"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	. "crype/api/gen_jet/crype_db/public/table"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CrypeServer struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

func (server *CrypeServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	paymentId := uuid.New()
	wallet, err := utils.GeneratePaymentAddress(req.Currency)
	if err != nil {
		return nil, err
	}
	stmt := PaymentAddresses.INSERT(PaymentAddresses.Address, PaymentAddresses.PrivateKey).VALUES(wallet.Address, wallet.PrivateKey)
	_, err = stmt.Exec(server.db)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:             fmt.Sprint(paymentId),
		PaymentAddress: fmt.Sprint(wallet.Address),
		OrderExpiration: timestamppb.New(
			time.Now().Add(time.Hour),
		),
		CreatedAt: timestamppb.Now(),
	}, nil
}

func main() {
	// The .env file is in the parent folder
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load .env file")
	}
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Initialize database connection, fails if not successful
	db := utils.ConnectDB()
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &CrypeServer{db: db})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
