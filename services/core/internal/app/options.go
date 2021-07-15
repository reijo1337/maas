package app

import "github.com/sirupsen/logrus"

type options struct {
	logger *logrus.Logger
}

type AppOption func(app *options)

func WithLogger(logger *logrus.Logger) AppOption {
	return func(o *options) {
		o.logger = logger
	}
}
