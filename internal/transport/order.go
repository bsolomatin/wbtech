package transport

import (
	"context"
	"dockertest/internal/models"
	"dockertest/internal/repository"
	"dockertest/pkg/logger"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

func (t *OrderTransport) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUid := r.URL.Query().Get("order_uid")
	if orderUid == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		t.logger.Error(r.Context(), "Fail to get order", zap.String("Order_uid", orderUid), zap.Error(fmt.Errorf("order_uid is empty")))
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
		t.logger.Error(r.Context(), "Fail to get order by order_uid", zap.String("Order_uid", orderUid), zap.Error(err))
		return
	}

	t.cache.Add(*order)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (t *OrderTransport) HomeTemplateHandler (w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		t.logger.Error(r.Context(), "Fail to parse template", zap.String("Template name", "index.html"), zap.Error(err))
		return 
	}

	temp.Execute(w, models.Order{})
}