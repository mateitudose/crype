package server

import (
	"context"
	"crype/server/config"
	"crype/server/service"
	"fmt"
	"log"
	"net"

	pb "crype/api/gen_proto"

	"google.golang.org/grpc"
)

type CrypeServer struct {
	pb.UnimplementedOrderServiceServer
	orderService service.OrderServiceInterface
}

func NewCrypeServer(orderService service.OrderServiceInterface) *CrypeServer {
	return &CrypeServer{orderService: orderService}
}

func (server *CrypeServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return server.orderService.CreateOrder(ctx, req)
}

func (server *CrypeServer) CheckOrderStatus(req *pb.CheckOrderStatusRequest, stream pb.OrderService_CheckOrderStatusServer) error {
	return server.orderService.CheckOrderStatus(req, stream)
}

func SetupServer(config *config.ServerConfig) error {
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", config.Port, err)
	}

	s := grpc.NewServer()

	orderRepository := service.NewOrderRepository(config.DB)
	orderService := service.NewOrderService(orderRepository)

	crypeServer := NewCrypeServer(orderService)
	pb.RegisterOrderServiceServer(s, crypeServer)

	log.Printf("Server listening on :%s", config.Port)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
