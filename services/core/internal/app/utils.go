package app

import "github.com/sirupsen/logrus"

func defaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	return logger
}
