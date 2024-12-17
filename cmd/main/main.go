package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"wbzerolevel/internal/config"
	"wbzerolevel/internal/repository"
	"wbzerolevel/internal/service"
	"wbzerolevel/internal/transport"
	"wbzerolevel/internal/transport/kafka"
	"wbzerolevel/internal/transport/middleware"
	"wbzerolevel/pkg/db"
	"wbzerolevel/pkg/logger"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New("WildberriesOrderService")
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	cfg, err := config.New()
	if cfg == nil {
		mainLogger.Error(ctx, "Fail to load config", zap.Error(err))
		panic("Fail to load config")
	}

	db, err := postgres.New(cfg.Config)
	if err != nil {
		mainLogger.Error(ctx, "Fail to connect to database", zap.Error(err))
		os.Exit(1)
	}
	orderRepository := repository.NewOrderRepository(db.Database)
	srv := service.NewOrderService(orderRepository, mainLogger)
	transp := transport.NewOrderTransport(srv, mainLogger)

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware(mainLogger))
	router.HandleFunc("/", transp.HomeTemplateHandler).Methods("GET")
	router.HandleFunc("/orders", transp.OrderTemplateHandler).Methods("GET")

	kafkaConsumer := kafka.NewReader(cfg.ConsumerConfig, srv.Process, mainLogger)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		kafkaConsumer.Consume(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), router)
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt)
	wg.Wait()
	<-graceSh

	mainLogger.Info(ctx, "Server stopped")
}
