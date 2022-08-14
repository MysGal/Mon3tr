package utils

import "go.uber.org/zap"

var GlobalLogger *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	suger := logger.Sugar()
	GlobalLogger = suger
}
