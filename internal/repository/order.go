package repository

import (
	"context"
	"dockertest/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	database *sqlx.DB
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
		return nil, err
	}
	defer tran.Rollback()

	const orderQuery = `INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, delivery_service, shard_key, sm_id, created_date, oof_shard)
              VALUES (:order_uid, :track_number, :entry, :locale, :customer_id, :delivery_service, :shard_key, :sm_id, :created_date, :oof_shard)`
	if _, err := tran.NamedExec(orderQuery, order); err != nil {
		fmt.Println("orders")
		fmt.Println(err)
		tran.Rollback()
		return nil, err
	}

	const deliveryQuery = `INSERT INTO delivery (order_id, name, phone, zip, city, address, region, email)
	VALUES (:order_id, :name, :phone, :zip, :city, :address, :region, :email)`
	order.Delivery.OrderUid = order.OrderUid
	if _, err := tran.NamedExec(deliveryQuery, order.Delivery); err != nil {
		fmt.Println("delivery")
		fmt.Println(err)
		tran.Rollback()
		return nil, err
	}

	const paymentQuery = `INSERT INTO payments (order_id, transaction, requestId, currency, provider, amount, paymentDateTime, bank, deliveryCost, goodsTotal, customFee)
	 VALUES (:order_id, :transaction, :requestId, :currency, :provider, :amount, :paymentDateTime, :bank, :deliveryCost, :goodsTotal, :customFee)`
	 order.Payment.OrderUid = order.OrderUid
	 if _, err := tran.NamedExec(paymentQuery, order.Payment); err != nil {
		fmt.Println("payments")
		fmt.Println(err)
		tran.Rollback()
		return nil, err
	 }

	const itemQuery = `INSERT INTO items (order_id, chrtId, trackNumber, price, rid, name, sale, size, totalPrice, nmId, brand, status)
	VALUES (:order_id, :chrtId, :trackNumber, :price, :rid, :name, :sale, :size, :totalprice, :nmId, :brand, :status)`
	for _, item := range order.Items {
		item.OrderUid = order.OrderUid
		if _, err := tran.NamedExec(itemQuery, item); err != nil {
			fmt.Println("items")
			fmt.Println(err)
			tran.Rollback();
			return nil, err
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
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
       

		if err := rows.StructScan(&order); err != nil {
			return nil, err
		}
        
        // if err = rows.StructScan(&item); err == nil {
        //     payload.Items = append(payload.Items, item)
        // }
		
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return &order, nil
}