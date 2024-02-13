package logger

import "go.uber.org/zap"

var log *zap.Logger

const defaultLogLevel = "INFO"

func NewLogger(level string) (*zap.Logger, error) {
	//lvl, err := zap.NewAtomicLevelAt(level)
	//if err != nil {
	//	return nil, err
	//}

	cfg := zap.NewProductionConfig()
	//cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	log = zl

	return zl, err
}

func Logger() *zap.Logger {
	if log != nil {
		return log
	}

	var err error
	log, err = NewLogger(defaultLogLevel)
	if err != nil {
		panic(err)
	}

	return log
}
