package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerKey = "logger"
	RequestId = "requestId"
	ServiceName = "service"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type logger struct {
	serviceName string
	logger *zap.Logger
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String(ServiceName, l.serviceName), zap.String(RequestId, ctx.Value(RequestId).(string)))
	l.logger.Info(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String(ServiceName, l.serviceName), zap.String(RequestId, ctx.Value(RequestId).(string)))
	l.logger.Error(msg, fields...)
}

func New(serviceName string) Logger {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	return &logger{
		serviceName : serviceName,
		logger: zapLogger,
	}
}

func GetLoggerFromCtx (ctx context.Context) Logger {
	return ctx.Value(LoggerKey).(Logger)

}