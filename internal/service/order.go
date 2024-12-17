package service

import (
	"context"
	"wbzerolevel/internal/models"
	"wbzerolevel/pkg/logger"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

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

func Validate(order models.Order) error {
	if err := ValidateOrder(order); err != nil {
		return fmt.Errorf("OrderService.ValidateOrder: %w", err)
	}

	if err := ValidatePayment(order.Payment); err != nil {
		return fmt.Errorf("OrderService.ValidatePayment: %w", err)
	}

	if err := ValidateDelivery(order.Delivery); err != nil {
		return fmt.Errorf("OrderService.ValidateDelivery: %w", err)
	}

	for _, item := range order.Items {
        if err := ValidateItem(item); err != nil {
            return fmt.Errorf("OrderService.ValidateItem: %w", err)
        }
    }

	return nil
}

func ValidateOrder(order models.Order) error {
	if order.OrderUid == "" {
        return fmt.Errorf("invalid order_uid")
    }

	if order.CustomerId == "" {
        return fmt.Errorf("empty customer_id")
    }
    if order.DeliveryService == "" {
        return fmt.Errorf("empty delivery_service")
    }
    if order.ShardKey == "" {
        return fmt.Errorf("empty shard key")
    }
	if order.CreatedDate.IsZero() || order.CreatedDate.After(time.Now()) {
        return fmt.Errorf("invalid created date")
    }

	return nil
}

func ValidateDelivery(delivery models.Delivery) error {
	return nil
}

func ValidatePayment(payment models.Payment) error {
	return nil
}

func ValidateItem(item models.Item) error {
	return nil
}

func isValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

func isValidCurrency(currency string) bool {
    validCurrencies := []string{"USD", "EUR", "RUB", "KZT", "CNY"}
    for _, c := range validCurrencies {
        if currency == c {
            return true
        }
    }
    return false
}

func isValidPhone(phone string) bool {
	var phoneRegex = regexp.MustCompile(`^\+\d{10,15}$`)
    return phoneRegex.MatchString(phone)
}

func isValidAmount(amount int) bool {
    return amount > 0
}

func isValidPaymentDT(paymentDT int64) bool {
    t := time.Unix(paymentDT, 0)
    return !t.IsZero() && t.After(time.Now())
}

func isValidCustomFee(customFee int) bool {
    return customFee >= 0
}

func isValidPrice(price int) bool {
    return price > 0
}

func isValidSale(sale int) bool {
    return sale >= 0 && sale <= 100
}

