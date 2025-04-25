package main

import (
	"context"
	"crype/api/gen_jet/crype_db/public/model"
	pb "crype/api/gen_proto"
	"crype/utils"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	// Dot import so that the Go code resembles native SQL
	. "crype/api/gen_jet/crype_db/public/table"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ServerConfig holds all configuration for the server
type ServerConfig struct {
	Port string
	DB   *sql.DB
}

// CrypeServer implements the OrderService gRPC server
type CrypeServer struct {
	pb.UnimplementedOrderServiceServer
	db *sql.DB
}

// NewCrypeServer creates a new instance of CrypeServer
func NewCrypeServer(db *sql.DB) *CrypeServer {
	return &CrypeServer{db: db}
}

// CreateOrder handles the creation of a new order
func (server *CrypeServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Generate a new payment ID
	paymentId := uuid.New()

	// Generate wallet address for the specified currency
	wallet, err := utils.GeneratePaymentAddress(req.Currency)
	if err != nil {
		log.Printf("Failed to generate payment address: %v", err)
		return nil, fmt.Errorf("failed to generate payment address: %w", err)
	}

	// Save payment address to database
	stmt := PaymentAddresses.INSERT(PaymentAddresses.Address, PaymentAddresses.PrivateKey).VALUES(wallet.Address, wallet.PrivateKey)
	_, err = stmt.Exec(server.db)
	if err != nil {
		log.Printf("Failed to insert payment address: %v", err)
		return nil, fmt.Errorf("failed to save payment address: %w", err)
	}

	// Set order creation and expiration times
	orderCreation := timestamppb.Now()
	orderExpiration := timestamppb.New(orderCreation.AsTime().Add(time.Hour))

	// Save order to database
	stmt = Orders.INSERT(Orders.ID, Orders.Amount, Orders.Currency, Orders.PaymentAddress, Orders.CreatedAt, Orders.OrderExpiration).VALUES(
		paymentId, req.Amount, req.Currency, wallet.Address, orderCreation.AsTime(), orderExpiration.AsTime())
	_, err = stmt.Exec(server.db)
	if err != nil {
		log.Printf("Failed to insert order: %v", err)
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Return order details
	return &pb.CreateOrderResponse{
		Id:              paymentId.String(),
		PaymentAddress:  fmt.Sprint(wallet.Address),
		CreatedAt:       orderCreation,
		OrderExpiration: orderExpiration,
	}, nil
}

// CheckOrderStatus streams status updates for a specific order
func (server *CrypeServer) CheckOrderStatus(req *pb.CheckOrderStatusRequest, stream pb.OrderService_CheckOrderStatusServer) error {
	// Parse and validate order ID
	orderId, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Invalid order ID format: %v", err)
		return fmt.Errorf("invalid order ID: %w", err)
	}

	// Query the order from database
	stmt := Orders.SELECT(Orders.AllColumns).WHERE(Orders.ID.EQ(postgres.UUID(orderId)))
	var orders []model.Orders
	err = stmt.Query(server.db, &orders)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return fmt.Errorf("database query error: %w", err)
	}

	// Verify the order exists
	if len(orders) == 0 {
		log.Printf("Order not found: %s", req.Id)
		return fmt.Errorf("order not found")
	}
	order := orders[0]

	// Convert order status to protobuf format
	orderStatus, ok := utils.GetProtobufOrderStatus(order.Status)
	if !ok {
		log.Printf("Cannot convert order status: %v", order.Status)
		return fmt.Errorf("cannot convert order status: %v", order.Status)
	}

	// Prepare and send initial response
	initialResponse := &pb.CheckOrderStatusResponse{
		Id:     order.ID.String(),
		Status: orderStatus,
		TxHash: order.TxHash,
	}

	// Send initial status of the order
	if err := stream.Send(initialResponse); err != nil {
		log.Printf("Failed to send initial status: %v", err)
		return err
	}

	// If the order is in a final state, end the stream
	if order.Status == model.OrderStatus_Confirmed ||
		order.Status == model.OrderStatus_Canceled ||
		order.Status == model.OrderStatus_Failed {
		return nil
	}

	// TODO: Add blockchain checking logic and stream the status changes
	// For now, we will simulate a status update after 5 seconds
	time.Sleep(5 * time.Second)
	// Simulate a status update
	order.Status = model.OrderStatus_Confirmed
	orderStatus, ok = utils.GetProtobufOrderStatus(order.Status)
	if !ok {
		log.Printf("Cannot convert order status: %v", order.Status)
		return fmt.Errorf("cannot convert order status: %v", order.Status)
	}
	txHash := "784c7f43639ad8c3951ad9a35890102c720937d21ce51c45498da053bb347160" // Mock transaction hash
	order.TxHash = &txHash
	// Prepare and send updated response
	updatedResponse := &pb.CheckOrderStatusResponse{
		Id:     order.ID.String(),
		Status: orderStatus,
		TxHash: order.TxHash,
	}
	if err := stream.Send(updatedResponse); err != nil {
		log.Printf("Failed to send updated status: %v", err)
	}

	return nil
}

// setupServer initializes and starts the gRPC server
func setupServer(config *ServerConfig) error {
	// Create network listener
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", config.Port, err)
	}

	// Create gRPC server
	s := grpc.NewServer()

	// Register our service implementation
	crypeServer := NewCrypeServer(config.DB)
	pb.RegisterOrderServiceServer(s, crypeServer)

	// Start serving requests
	log.Printf("Server listening on :%s", config.Port)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Create server configuration
	config := &ServerConfig{
		Port: os.Getenv("CRYPE_PORT"), 
		DB:   utils.ConnectDB(),
	}

	// Start the server
	if err := setupServer(config); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
