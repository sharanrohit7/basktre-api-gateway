package logger

import (
	"context"
	"sync"

	"github.com/basktre/api-gateway/pkg/requestid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once   sync.Once
	logger *zap.Logger
)

func Initialize(level, format string) error {
	var initErr error
	once.Do(func() {
		lvl := zapcore.InfoLevel
		if err := lvl.Set(level); err != nil {
			lvl = zapcore.InfoLevel
		}

		var cfg zap.Config
		if format == "json" {
			cfg = zap.NewProductionConfig()
		} else {
			cfg = zap.NewDevelopmentConfig()
		}
		cfg.Level = zap.NewAtomicLevelAt(lvl)

		l, err := cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			initErr = err
			return
		}
		logger = l
	})
	return initErr
}

func L() *zap.Logger {
	if logger == nil {
		return zap.NewNop()
	}
	return logger
}

func Info(msg string, fields ...zap.Field)  { L().Info(msg, fields...) }
func Error(msg string, fields ...zap.Field) { L().Error(msg, fields...) }
func Fatalf(format string, args ...interface{}) {
	L().Sugar().Fatalf(format, args...)
}
func Infof(format string, args ...interface{}) { L().Sugar().Infof(format, args...) }
func Warnf(format string, args ...interface{}) { L().Sugar().Warnf(format, args...) }

func WithContext(ctx context.Context) *zap.Logger {
	id := requestid.GetRequestID(ctx)
	if id == "" {
		return L()
	}
	return L().With(zap.String("request_id", id))
}

func InfofWithContext(ctx context.Context, format string, args ...interface{}) {
	WithContext(ctx).Sugar().Infof(format, args...)
}
