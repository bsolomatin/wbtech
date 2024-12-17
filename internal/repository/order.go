package repository

import (
	"context"
	"wbzerolevel/internal/models"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	database *sqlx.DB
}

type OrderCacheEntry struct {
	Order models.Order
	timer *time.Timer
}

type OrderCache struct {
	mu     *sync.RWMutex
	Orders map[string]*OrderCacheEntry
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		Orders: make(map[string]*OrderCacheEntry),
		mu:     &sync.RWMutex{},
	}
}

func (o *OrderCache) Add(data models.Order, ttl ...time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	ttlDuration := 12 * time.Hour
	if ttl != nil {
		ttlDuration = ttl[0]
	}

	if order, exists := o.Orders[data.OrderUid]; exists {
		order.timer.Stop()
	}

	timer := time.AfterFunc(ttlDuration, func() {
		o.Invalidate(data.OrderUid)
	})

	o.Orders[data.OrderUid] = &OrderCacheEntry{
		Order: data,
		timer: timer,
	}
}

func (o *OrderCache) Get(orderUid string) (models.Order, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	order, exists := o.Orders[orderUid]
	if !exists {
		return models.Order{}, false
	}

	return order.Order, true
}

func (o *OrderCache) Invalidate(orderUid string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if order, exists := o.Orders[orderUid]; exists {
		order.timer.Stop()
		delete(o.Orders, orderUid)
	}
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{
		database: db,
	}
}

func (s *OrderRepository) CreateNewOrder(ctx context.Context, order models.Order) (*models.Order, error) {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("OrderRepository.CreateNewOrder: %w", err)
	}

	const query = `INSERT INTO orders (order_uid, data) VALUES ($1, $2)`
	_, err = s.database.ExecContext(ctx, query, order.OrderUid, jsonData)
	if err != nil {
		return nil, fmt.Errorf("OrderRepository.CreateNewOrder: %w", err)
	}

	return &order, nil
}

func (s *OrderRepository) FindByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	var order models.Order
	var jsonData json.RawMessage
    const query = "SELECT data FROM orders WHERE order_uid = $1"
    err := s.database.GetContext(ctx, &jsonData, query, orderUid)
    if err != nil {
        return nil, fmt.Errorf("OrderRepository.FindByUid: %w", err)
    }

    if err = json.Unmarshal(jsonData, &order); err != nil {
        return nil, fmt.Errorf("OrderRepository.FindByUid: %w", err)
	}

    return &order, nil
}
