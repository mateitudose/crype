package service

import (
	"context"
	"crype/api/gen_jet/crype_db/public/model"
	pb "crype/api/gen_proto"
	"crype/server/utils"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService struct {
	repository OrderRepositoryInterface
}

func NewOrderService(repository OrderRepositoryInterface) *OrderService {
	return &OrderService{repository: repository}
}

// CreateOrder handles the creation of a new order
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Generate a new payment ID
	paymentId := uuid.New()

	// Generate wallet address for the specified currency
	wallet, err := utils.GeneratePaymentAddress(req.Currency)
	if err != nil {
		log.Printf("Failed to generate payment address: %v", err)
		return nil, fmt.Errorf("failed to generate payment address: %w", err)
	}

	// Save payment address to database
	err = s.repository.SavePaymentAddress(wallet.Address, wallet.PrivateKey)
	if err != nil {
		log.Printf("Failed to insert payment address: %v", err)
		return nil, err
	}

	// Set order creation and expiration times
	orderCreation := timestamppb.Now()
	orderExpiration := timestamppb.New(orderCreation.AsTime().Add(time.Hour))

	// Save order to database
	err = s.repository.SaveOrder(
		paymentId,
		req.Amount,
		req.Currency,
		wallet.Address,
		orderCreation.AsTime(),
		orderExpiration.AsTime(),
	)
	if err != nil {
		log.Printf("Failed to insert order: %v", err)
		return nil, err
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
func (s *OrderService) CheckOrderStatus(req *pb.CheckOrderStatusRequest, stream pb.OrderService_CheckOrderStatusServer) error {
	// Parse and validate order ID
	orderId, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("Invalid order ID format: %v", err)
		return fmt.Errorf("invalid order ID: %w", err)
	}

	// Query the order from database
	order, err := s.repository.GetOrderByID(orderId)
	if err != nil {
		log.Printf("Failed to get order: %v", err)
		return err
	}

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
	txHash := "784c7f43639ad8c3951ad9a35890102c720937d21ce51c45498da053bb347160" // Mock transaction hash
	order.TxHash = &txHash

	// Update the order in the database
	err = s.repository.UpdateOrderStatus(order.ID, order.Status, order.TxHash)
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	orderStatus, ok = utils.GetProtobufOrderStatus(order.Status)
	if !ok {
		log.Printf("Cannot convert order status: %v", order.Status)
		return fmt.Errorf("cannot convert order status: %v", order.Status)
	}

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
