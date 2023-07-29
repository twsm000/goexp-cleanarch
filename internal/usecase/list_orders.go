package usecase

import (
	"github.com/twsm000/goexp-cleanarch/internal/entity"
)

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func (uc *ListOrdersUseCase) Execute() ([]*OrderOutputDTO, error) {
	orders, err := uc.OrderRepository.ListAll()
	if err != nil {
		return nil, err
	}

	result := make([]*OrderOutputDTO, len(orders))
	for i := 0; i < len(orders); i++ {
		result[i] = &OrderOutputDTO{
			ID:         orders[i].ID,
			Price:      orders[i].Price,
			Tax:        orders[i].Tax,
			FinalPrice: orders[i].FinalPrice,
		}
	}

	return result, nil
}
