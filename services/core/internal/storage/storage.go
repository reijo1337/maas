package storage

import (
	"context"
	"io"
)

type Manager interface {
	io.Closer
	Ping(ctx context.Context) error
	RegisterClient(ctx context.Context, host, name, key string) error
}

func New() (Manager, error) {
	return newPostgresManager()
}
