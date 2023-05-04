package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.Logger
type SugaredLogger = zap.SugaredLogger
type Option = zap.Option

func NewLogger(options ...Option) (*Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderCfg
	return cfg.Build()
}

func NewSugaredLogger(options ...Option) (*SugaredLogger, error) {
	logger, err := NewLogger(options...)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
