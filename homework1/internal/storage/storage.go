package storage

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/alextonkonogov/gb-go-observability/homework1/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDBConn(ctx context.Context, appConfig *config.AppConfig) (dbpool *pgxpool.Pool, err error) {
	url := "postgres://postgres:password@localhost:5432/postgres"

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return
	}

	cfg.MaxConns = int32(appConfig.MaxConns)
	cfg.MinConns = int32(appConfig.MinConns)
	cfg.HealthCheckPeriod = 1 * time.Minute
	cfg.MaxConnLifetime = 24 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		Timeout:   cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	dbpool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return
	}

	return
}
