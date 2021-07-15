package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/kelseyhightower/envconfig"
)

type postgres struct {
	db *pgx.Conn
}

type postgresEnvConfig struct {
	Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     uint16 `envconfig:"POSTGRES_PORT" default:"5432"`
	Name     string `envconfig:"POSTGRES_DB_NAME" defautl:"julia"`
	User     string `envconfig:"POSTGRES_DB_USER" required:"true"`
	Password string `envconfig:"POSTGRES_DB_PASSWORD" required:"true"`
}

func parseEnvConfig() (*postgresEnvConfig, error) {
	out := &postgresEnvConfig{}
	if err := envconfig.Process("", out); err != nil {
		envconfig.Usage("", out)
		return nil, err
	}

	return out, nil
}

func newPostgresManager() (Manager, error) {
	cfg, err := parseEnvConfig()
	if err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}

	db, err := pgx.Connect(pgx.ConnConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Database: cfg.Name,
		User:     cfg.User,
		Password: cfg.Password,
	})

	if err != nil {
		return nil, fmt.Errorf("open pgx connect: %v", err)
	}

	return &postgres{db: db}, nil
}

func (p *postgres) Close() error {
	return p.db.Close()
}

func (p *postgres) Ping(ctx context.Context) error {
	return p.db.Ping(ctx)
}

func (p *postgres) RegisterClient(ctx context.Context, host, name, key string) error {
	md5key := getMD5Hash(key)
	i := 0
	if err := p.db.QueryRow("select 1 from clients where name=$1 and key=$2", name, md5key).Scan(&i); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("scan existing clients: %v", err)
		}
	} else {
		return nil
	}

	if _, err := p.db.Exec("insert into clients (host, name, key) values ($1, $2, $3)", host, name, md5key); err != nil {
		return fmt.Errorf("inseret new client: %v", err)
	}

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte("md5" + text))
	return hex.EncodeToString(hash[:])
}
