package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.Logger
type SugaredLogger = zap.SugaredLogger
type Option func(cfg *zap.Config)

type Level string

var (
	DEBUG Level = "debug"
	INFO  Level = "info"
	ERROR Level = "error"
)

func WithLogLevel(lvl Level) Option {
	return func(cfg *zap.Config) {
		switch lvl {
		case DEBUG:
			cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		case ERROR:
			cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		default:
			cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		}
	}
}

func NewLogger(options ...Option) (*Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderCfg
	for _, applyOpt := range options {
		applyOpt(&cfg)
	}

	return cfg.Build()
}

func NewSugaredLogger(options ...Option) (*SugaredLogger, error) {
	logger, err := NewLogger(options...)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
