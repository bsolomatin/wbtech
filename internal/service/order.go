package service

import (
	"context"
	"dockertest/internal/models"
	"dockertest/pkg/logger"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
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
	// if err := s.Validate(ctx, order); err != nil {
	// 	return nil, fmt.Errorf("OrderService.CreateNewOrder: %s", err)
	// }
	return s.Repo.CreateNewOrder(ctx, order)
}

// func (s *OrderService) Validate(ctx context.Context, order models.Order) error {
// 	return nil
// }

func (s *OrderService) FindByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	return s.Repo.FindByUid(ctx, orderUid)
}

func (s *OrderService) Process(ctx context.Context, message kafka.Message) error {
	data := models.Order{}
	if err := json.Unmarshal(message.Value, &data); err != nil {
		return fmt.Errorf("OrderService.Process: %s", err)
	}
	_, err := s.CreateNewOrder(ctx, data)
	if err != nil {
		return fmt.Errorf("OrderService.Process: %s", err)
	}
	return nil 
}