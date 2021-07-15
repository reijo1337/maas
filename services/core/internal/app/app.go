package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Application struct {
	cfg        *config
	logger     *logrus.Logger
	server     *rpc.Server
	httpServer *http.Server
}

type config struct {
	Port            string        `envconfig:"APPLICATION_PORT" default:":1337"`
	ShutdownTimeout time.Duration `envconfig:"APPLICATION_SHUTDOWN_TIMEOUT" default:"5s"`
}

func parseEnvConfig() (*config, error) {
	out := &config{}
	if err := envconfig.Process("", out); err != nil {
		envconfig.Usage("", out)
		return nil, err
	}

	return out, nil
}

func New(opts ...AppOption) (*Application, error) {
	cfg, err := parseEnvConfig()
	if err != nil {
		return nil, fmt.Errorf("parse env config: %v", err)
	}

	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	if o.logger == nil {
		o.logger = defaultLogger()
	}

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	return &Application{
		cfg:    cfg,
		logger: o.logger,
		server: s,
	}, nil
}

func (a *Application) RegisterRPC(rcsvr interface{}, name string) error {
	return a.server.RegisterService(rcsvr, name)
}

func (a *Application) Run() {
	r := mux.NewRouter()
	r.Handle("/rpc", a.server)

	a.httpServer = &http.Server{
		Addr:    a.cfg.Port,
		Handler: r,
	}

	go func() {
		a.logger.Infof("start listening on port %s", a.cfg.Port)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.WithError(err).Error("listen port %s", a.cfg.Port)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	a.logger.Info("termination signal received")

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.WithError(err).Error("shutdown server")
	}
}
