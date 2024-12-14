package service

import (
	"context"
	"dockertest/internal/models"
	"dockertest/pkg/logger"
)

type OrderRepository interface {
	CreateNewOrder(ctx context.Context, order models.Order) (*models.Order, error)
	FindByUid(ctx context.Context, orderUid string) (*models.Order, error)
}

type OrderService struct {
	Logger logger.Logger
	Repo OrderRepository
}

func NewOrderService(repo OrderRepository, logger logger.Logger) *OrderService {
	return &OrderService{
		Repo: repo,
		Logger: logger,
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