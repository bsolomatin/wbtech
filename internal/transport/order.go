package transport

import (
	"context"
	"dockertest/internal/models"
	"dockertest/internal/repository"
	"dockertest/pkg/logger"
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
	logger  logger.Logger
	cache   *repository.OrderCache
}

func NewOrderTransport(service OrderService, log logger.Logger) *OrderTransport {
	return &OrderTransport{
		service: service,
		logger:  log,
		cache:   repository.NewOrderCache(),
	}
}

func (t *OrderTransport) GetOrder(ctx context.Context, orderUid string) (*models.Order, error) {
	if order, exists := t.cache.Get(orderUid); exists {
		return &order, nil
	}

	order, err := t.service.FindByUid(ctx, orderUid)
	if err != nil {
		return nil, err
	}

	t.cache.Add(*order)

	return order, nil
}

func (t *OrderTransport) HomeTemplateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		t.logger.Error(r.Context(), "Fail to parse template", zap.String("Template name", "index.html"), zap.Error(err))
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		t.logger.Error(r.Context(), "Fail to execute template", zap.String("Template name", "index.html"), zap.Error(err))
	}
}

func (t *OrderTransport) OrderTemplateHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/order.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.logger.Error(r.Context(), "Fail to parse template", zap.String("Template name", "order.html"), zap.Error(err))
		return
	}

	orderUid := r.URL.Query().Get("order_uid")
	if orderUid == "" {
		w.WriteHeader(http.StatusBadRequest)
		tmpl.Execute(w, models.PageData{Error: err.Error()})
		t.logger.Error(r.Context(), "Fail to get order", zap.String("Order_uid", orderUid), zap.Error(fmt.Errorf("order_uid is empty")))
		return
	}
	order, err := t.GetOrder(r.Context(), orderUid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmpl.Execute(w, models.PageData{Error: err.Error()})
		t.logger.Error(r.Context(), "Fail to get order", zap.String("Order_uid", orderUid), zap.Error(err))
		return
	}
	if err := tmpl.Execute(w, models.PageData{Order: *order}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		t.logger.Error(r.Context(), "Fail to execute template", zap.String("Template name", "order.html"), zap.Error(err))
	}
}
