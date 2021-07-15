package main

import (
	"github.com/reijo1337/maas/services/core/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := defaultLogger()
	a, err := app.New(app.WithLogger(logger))

	if err != nil {
		logger.WithError(err).Fatal("init application")
	}

	a.Run()
}

func defaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	return logger
}
