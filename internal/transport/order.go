package transport

import (
	"context"
	"dockertest/internal/models"
	"dockertest/internal/repository"
	"dockertest/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type OrderService interface {
	CreateNewOrder(ctx context.Context, order models.Order) (*models.Order, error)
	FindByUid(ctx context.Context, orderUid string) (*models.Order, error)
}

type OrderTransport struct {
	service OrderService
	logger logger.Logger
	cache *repository.OrderCache
}

func NewOrderTransport(service OrderService, log logger.Logger) *OrderTransport {
    return &OrderTransport{
		service: service,
		logger: log,
		cache: repository.NewOrderCache(),
	}
}

func (t *OrderTransport) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var payload models.Order
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		t.logger.Error(r.Context(), "Fail to decode request", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payload)
}

func (t *OrderTransport) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUid := r.URL.Query().Get("order_uid")
	if orderUid == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		t.logger.Error(r.Context(), "Fail to get order", zap.String("Order_uid", orderUid))
		return
	}

	if order, exists := t.cache.Get(orderUid); exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)
		return
	}

	order, err := t.service.FindByUid(r.Context(), orderUid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		t.logger.Error(r.Context(), fmt.Sprintf("Fail to get order by order_uid %s", orderUid), zap.Error(err))
		return
	}

	t.cache.Add(*order, 15 * time.Second)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func Validate () error{
	return nil
}