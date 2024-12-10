package service

import (
	"context"
	"dockertest/internal/models"
)

type OrderRepository interface {
	CreateNewOrder(ctx context.Context, order models.Order) (*models.Order, error)
	FindByUid(ctx context.Context, orderUid string) (*models.Order, error)
}

type OrderService struct {
	Repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		Repo: repo,
	}
}

func (s *OrderService) CreateNewOrder(ctx context.Context, order models.Order) (*models.Order, error){
	if err := s.Validate(ctx, order); err != nil {
		return nil, err
	}
	return s.Repo.CreateNewOrder(ctx, order)
}

func (s *OrderService) Validate(ctx context.Context, order models.Order) error {
	return nil
}

func (s *OrderService) FindByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	return s.Repo.FindByUid(ctx, orderUid)
}