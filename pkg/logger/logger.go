package logger

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	var logger *zap.Logger
	var err error

	if os.Getenv("APP_ENV") == "dev" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	Log = logger.Sugar()
}

func Sync() {
	if Log != nil {
		Log.Sync()
	}
}
