package service

import (
	"context"
	"crype/api/gen_jet/crype_db/public/model"
	pb "crype/api/gen_proto"

	"github.com/google/uuid"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	CheckOrderStatus(req *pb.CheckOrderStatusRequest, stream pb.OrderService_CheckOrderStatusServer) error
}

type OrderRepositoryInterface interface {
	SavePaymentAddress(address, privateKey string) error
	SaveOrder(id uuid.UUID, amount float64, currency, paymentAddress string, createdAt, expiresAt interface{}) error
	GetOrderByID(id uuid.UUID) (*model.Orders, error)
	UpdateOrderStatus(id uuid.UUID, status model.OrderStatus, txHash *string) error
}
