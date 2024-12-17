package main

import (
	"context"
	"dockertest/internal/config"
	"dockertest/internal/repository"
	"dockertest/internal/service"
	"dockertest/internal/transport"
	"dockertest/internal/transport/kafka"
	"dockertest/pkg/db"
	"dockertest/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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
	router.Use(loggingMiddleware(mainLogger))
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

func loggingMiddleware(logger logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info(r.Context(), "Request received", zap.String("Method", r.Method), zap.String("URL", r.URL.Path))
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Info(r.Context(), "Request processed", zap.String("Method", r.Method), zap.String("URL", r.URL.Path), zap.String("Processed time", duration.String()))
		})
	}
}
