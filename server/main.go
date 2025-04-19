package main

import (
	"context"
	"crype/api/gen_jet/crype_db/public/model"
	pb "crype/api/gen_proto"
	"crype/utils"
	"database/sql"
	"fmt"
	"github.com/go-jet/jet/v2/postgres"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"

	// Dot import so that the Go code resembles native SQL
	. "crype/api/gen_jet/crype_db/public/table"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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
	orderCreation := timestamppb.Now()
	orderExpiration := timestamppb.New(orderCreation.AsTime().Add(time.Hour))
	stmt = Orders.INSERT(Orders.ID, Orders.Amount, Orders.Currency, Orders.PaymentAddress, Orders.CreatedAt, Orders.OrderExpiration).VALUES(
		paymentId, req.Amount, req.Currency, wallet.Address, orderCreation.AsTime(), orderExpiration.AsTime())
	_, err = stmt.Exec(server.db)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:              paymentId.String(),
		PaymentAddress:  fmt.Sprint(wallet.Address),
		CreatedAt:       orderCreation,
		OrderExpiration: orderExpiration,
	}, nil
}

func (server *CrypeServer) CheckOrderStatus(req *pb.CheckOrderStatusRequest, stream pb.OrderService_CheckOrderStatusServer) error {
	orderId, err := uuid.Parse(req.Id)
	if err != nil {
		return fmt.Errorf("invalid order ID: %s", err)
	}
	stmt := Orders.SELECT(Orders.AllColumns).WHERE(Orders.ID.EQ(postgres.UUID(orderId)))
	var orders []model.Orders
	err = stmt.Query(server.db, &orders)
	if err != nil {
		return fmt.Errorf("database query error: %s", err)
	}
	if len(orders) == 0 {
		return fmt.Errorf("order not found")
	}
	order := orders[0]
	// Safely map the status enum from SQL to the corresponding protobuf enum
	orderStatusMap := map[model.OrderStatus]pb.OrderStatus{
		model.OrderStatus_Pending:    pb.OrderStatus_PENDING,
		model.OrderStatus_Processing: pb.OrderStatus_PROCESSING,
		model.OrderStatus_Confirmed:  pb.OrderStatus_CONFIRMED,
		model.OrderStatus_Failed:     pb.OrderStatus_FAILED,
		model.OrderStatus_Canceled:   pb.OrderStatus_CANCELED,
	}
	orderStatus, ok := orderStatusMap[order.Status]
	if !ok {
		return fmt.Errorf("cannot convert order status: %v", order.Status)
	}
	initialResponse := &pb.CheckOrderStatusResponse{
		Id:     order.ID.String(),
		Status: orderStatus,
		TxHash: order.TxHash,
	}
	// Send initial status of the order
	err = stream.Send(initialResponse)
	if err != nil {
		return err
	}
	// If the order is confirmed, canceled or failed, end the stream
	if order.Status == model.OrderStatus_Confirmed ||
		order.Status == model.OrderStatus_Canceled ||
		order.Status == model.OrderStatus_Failed {
		return nil
	}
	// TODO: Add blockchain checking logic and stream the status changes
	return nil
}

func main() {
	// This loads the environment variables from the .env file located in the root directory
	err := godotenv.Load()
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
