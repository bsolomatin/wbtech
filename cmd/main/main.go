package main

import (
	"context"
	"dockertest/internal/config"
	"dockertest/internal/models"
	"dockertest/internal/repository"
	"dockertest/internal/service"
	"dockertest/pkg/db"
	"dockertest/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/rand"
)

func messageHandler(w http.ResponseWriter, rq *http.Request) {
	messageId, err := strconv.Atoi(rq.URL.Query().Get("messageId"))
	if err != nil || messageId < 1 {
		fmt.Printf("err %s", err)
		http.NotFound(w, rq)
		return
	}
	fmt.Println(messageId)
}

func main() {
	fmt.Println("POEHALI")
	//comeback later
	ctx := context.Background()
	mainLogger := logger.New("WilberriesOrderService")
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	cfg := config.New()
	if cfg == nil {
		panic("Fail ot load config")
	}
	//comeback later

	// router := mux.NewRouter()
	// logger := setupLogger()
	// router.Use(loggingMiddleware(logger))
	// fmt.Println("Hello, world")
	// router.HandleFunc("/", messageHandler)
	// http.ListenAndServe(":8080", router)
	// cfg := config.New()
	// if cfg == nil {
	// 	panic("Fail to load config")
	// }


	//comeback later
	db, err := postgres.New(cfg.Config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	orderRepository := repository.NewOrderRepository(db.Database)
	srv := service.NewOrderService(orderRepository)
	srv.CreateNewOrder(ctx, generateRandomOrder())
	//comeback later
	//generateRandomOrder()
	// type Hello struct {
	// 	Field int `db:hello`
	// }
}

func setupLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Ошибка конфигурирования логгера: %s\n", err)
	}

	return logger
}

func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info("Получен запрос", zap.String("Метод", r.Method), zap.String("URL", r.URL.Path))
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Info("Обработан запрос", zap.String("Метод", r.Method), zap.String("URL", r.URL.Path), zap.String("Длительность обработки", duration.String()))
		})
	}
}

func generateRandomOrder() models.Order {
	now := time.Now()

	orderUID := uuid.New().String()

	// Генерируем случайный трек-номер
	trackNumber := fmt.Sprintf("TRACK%d", rand.Intn(1000000))

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
