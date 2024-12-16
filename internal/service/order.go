package service

import (
	"context"
	"dockertest/internal/models"
	"dockertest/pkg/logger"
	"encoding/json"
	"fmt"
	"regexp"

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
	return s.Repo.CreateNewOrder(ctx, order)
}

func (s *OrderService) FindByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	return s.Repo.FindByUid(ctx, orderUid)
}

func (s *OrderService) Process(ctx context.Context, message kafka.Message) error {
	data := models.Order{}
	if err := json.Unmarshal(message.Value, &data); err != nil {
		return fmt.Errorf("OrderService.Process: %w", err)
	}
	if err := Validate(data); err != nil {
		return fmt.Errorf("OrderService.Process: %w", err)
	}
	_, err := s.CreateNewOrder(ctx, data)
	if err != nil {
		return fmt.Errorf("OrderService.Process: %w", err)
	}
	
	return nil 
}

func Validate (order models.Order) error {
	if order.OrderUid == "" {
		return fmt.Errorf("orderUid is empty")
	}
	if !isValidEmail(order.Delivery.Email) {
		return fmt.Errorf("invalid email")
	}

	return nil
}

func isValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}