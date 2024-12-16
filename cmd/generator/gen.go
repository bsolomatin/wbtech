package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"dockertest/internal/models"
	"math/rand"
	"time"
)

func generateRandomOrder() models.Order {
	now := time.Now()
	orderUID := uuid.New().String()
	trackNumber := fmt.Sprintf("TRACK%d", rand.Intn(1000000))
	seed := time.Now().UnixNano()
    rand := rand.New(rand.NewSource(seed))

	return models.Order{
		OrderUid:    orderUID,
		TrackNumber: trackNumber,
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    fmt.Sprintf("Customer %d", rand.Intn(1000)),
			Phone:   fmt.Sprintf("+7%d", rand.Intn(999999999)),
			Zip:     fmt.Sprintf("%d", rand.Intn(999999)),
			City:    "Moscow",
			Address: fmt.Sprintf("Street %d", rand.Intn(100)),
			Region:  "Moscow Region",
			Email:   fmt.Sprintf("customer%d@example.com", rand.Intn(1000)),
		},
		Payment: models.Payment{
			Transaction:  orderUID,
			RequestId:    "",
			Currency:     "RUB",
			Provider:     "wbpay",
			Amount:       float64(rand.Intn(10000)),
			PaymentDateTime:    now.Unix(),
			Bank:         "alpha",
			DeliveryCost: float64(rand.Intn(1000)),
			GoodsTotal:   float64(rand.Intn(9000)),
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtId:      rand.Intn(9999999),
				TrackNumber: trackNumber,
				Price:       float64(rand.Intn(1000)),
				Rid:         uuid.New().String(),
				Name:        fmt.Sprintf("Product %d", rand.Intn(100)),
				Sale:        rand.Intn(90),
				Size:        "0",
				TotalPrice:  float64(rand.Intn(900)),
				NmId:        rand.Intn(9999999),
				Brand:       fmt.Sprintf("Brand %d", rand.Intn(10)),
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        fmt.Sprintf("customer_%d", rand.Intn(1000)),
		DeliveryService:   "meest",
		ShardKey:          fmt.Sprintf("%d", rand.Intn(10)),
		SmId:              99,
		CreatedDate:       now,
		OofShard:          "1",
	}
}

func main() {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "LEVELZERO",
	})
	defer w.Close()
	ctx := context.Background()
	for {
		order := generateRandomOrder()

		orderJSON, err := json.Marshal(order)
		if err != nil {
			fmt.Printf("Error marshaling order: %w", err)
			continue
		}

		err = w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(order.OrderUid),
			Value: orderJSON,
		})

		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		} else {
			fmt.Printf("Sent order: %s\n", order.OrderUid)
		}

		time.Sleep(time.Second)
	}
}
