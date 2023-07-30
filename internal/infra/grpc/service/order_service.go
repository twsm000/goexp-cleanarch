package service

import (
	"context"

	"github.com/twsm000/goexp-cleanarch/internal/infra/grpc/pb"
	"github.com/twsm000/goexp-cleanarch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewOrderService(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, _ *pb.EmptyMessageRequest) (*pb.OrderListResponse, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	resp := pb.OrderListResponse{
		Orders: make([]*pb.Order, len(orders)),
	}

	for i := 0; i < len(orders); i++ {
		resp.Orders[i] = &pb.Order{
			Id:         orders[i].ID,
			Price:      orders[i].Price,
			Tax:        orders[i].Tax,
			FinalPrice: orders[i].FinalPrice,
		}
	}

	return &resp, nil
}
