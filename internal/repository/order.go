package repository

import (
	"context"
	"dockertest/internal/models"
	"fmt"
	"sync"
	"time"
	"github.com/jmoiron/sqlx"
	"math/rand"
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

func (o *OrderCache) Add(data models.Order, ttl time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if order, exists := o.Orders[data.OrderUid]; exists {
		order.timer.Stop()
	} else {
		jitter := time.Duration(rand.Int63n(int64(ttl) / 10)) // Jitter up to 10% of TTL //todo 
        adjustedTTL := ttl + jitter
		timer := time.AfterFunc(adjustedTTL, func() {
			o.Invalidate(data.OrderUid)

		})
		o.Orders[data.OrderUid] = &OrderCacheEntry{
			Order: data,
			timer: timer,
		}
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
	//jsonItems, err := json.Marshal(order.i)
	tran, err := s.database.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("OrderRepository.CreateNewOrder: %s", err)
	}
	defer tran.Rollback()

	const orderQuery = `INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, delivery_service, shard_key, sm_id, created_date, oof_shard)
              VALUES (:order_uid, :track_number, :entry, :locale, :customer_id, :delivery_service, :shard_key, :sm_id, :created_date, :oof_shard)`
	if _, err := tran.NamedExec(orderQuery, order); err != nil {
		tran.Rollback()
		return nil, fmt.Errorf("OrderRepository.CreateNewOrder: %s", err)
	}

	const deliveryQuery = `INSERT INTO delivery (order_id, name, phone, zip, city, address, region, email)
	VALUES (:order_id, :name, :phone, :zip, :city, :address, :region, :email)`
	order.Delivery.OrderUid = order.OrderUid
	if _, err := tran.NamedExec(deliveryQuery, order.Delivery); err != nil {
		tran.Rollback()
		return nil, fmt.Errorf("OrderRepository.CreateNewOrder: %s", err)
	}

	const paymentQuery = `INSERT INTO payments (order_id, transaction, requestId, currency, provider, amount, paymentDateTime, bank, deliveryCost, goodsTotal, customFee)
	 VALUES (:order_id, :transaction, :requestId, :currency, :provider, :amount, :paymentDateTime, :bank, :deliveryCost, :goodsTotal, :customFee)`
	order.Payment.OrderUid = order.OrderUid
	if _, err := tran.NamedExec(paymentQuery, order.Payment); err != nil {
		tran.Rollback()
		return nil, fmt.Errorf("OrderRepository.CreateNewOrder: %s", err)
	}

	const itemQuery = `INSERT INTO items (order_id, chrtId, trackNumber, price, rid, name, sale, size, totalPrice, nmId, brand, status)
	VALUES (:order_id, :chrtId, :trackNumber, :price, :rid, :name, :sale, :size, :totalprice, :nmId, :brand, :status)`
	for _, item := range order.Items {
		item.OrderUid = order.OrderUid
		if _, err := tran.NamedExec(itemQuery, item); err != nil {
			tran.Rollback()
			return nil, fmt.Errorf("OrderRepository.CreateNewOrder: %s", err)
		}
	}

	tran.Commit()
	return &order, nil
}

func (s *OrderRepository) FindByUid(ctx context.Context, orderUid string) (*models.Order, error) {
	var order models.Order
	const query = "SELECT order.*, delivery.*, payment.*, item.* FROM order INNER JOIN delivery ON order.orderUid = delivery.orderUid"
	rows, err := s.database.QueryxContext(ctx, query, orderUid)
	if err != nil {
		return nil, fmt.Errorf("OrderRepository.FindByUid: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(&order); err != nil {
			return nil, fmt.Errorf("OrderRepository.FindByUid: %s", err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("OrderRepository.FindByUid: %s", err)
	}

	return &order, nil
}
