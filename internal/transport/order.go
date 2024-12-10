package transport

import (
	"dockertest/internal/models"
	"dockertest/internal/service"
	"encoding/json"
	"net/http"
)

type OrderTransport struct {
	service *service.OrderService
}

func NewOrderTransport(service *service.OrderService) *OrderTransport {
    return &OrderTransport{service: service}
}

func (t *OrderTransport) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var payload models.Order
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payload)
}

func (t *OrderTransport) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUid := r.URL.Query().Get("order_uid")
	order, err := t.service.FindByUid(r.Context(), orderUid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}