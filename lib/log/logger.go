package log

import "go.uber.org/zap"

type Logger = zap.Logger
type SugaredLogger = zap.SugaredLogger
type Option = zap.Option

func NewLogger(options ...Option) (*Logger, error) {
	return zap.NewProduction(options...)
}

func NewSugaredLogger(options ...Option) (*SugaredLogger, error) {
	logger, err := NewLogger(options...)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
