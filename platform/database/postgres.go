package database

import (
	"context"
	"fmt"
	"time"

	"zocket/pkg/configs"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ConnectPostgres func for connection to PostgreSQL database.
func ConnectPostgres(ctx context.Context, config *configs.PostgresConfig) (*pgxpool.Pool, error) {

	// Build PostgreSQL connection URL.
	postgresConnURL, err := connectionURLBuilder(config)
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgreSQL.
	pool, err := pgxpool.Connect(ctx, postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Try to ping database.
	if err := pool.Ping(ctx); err != nil {
		defer pool.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return pool, nil
}
func connectionURLBuilder(config *configs.PostgresConfig) (string, error) {

	url := fmt.Sprintf(
		"user=%s password=%s host=%s port=%v dbname=%s sslmode=%s TimeZone=%s pool_max_conns=%v pool_max_conn_idle_time=%v pool_max_conn_lifetime=%v",
		config.UserName,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
		config.TimeZone,
		config.MaxConnections,
		time.Duration(time.Duration(config.MaxIdleConnectionTimeoutMinutes)*time.Minute),
		time.Duration(time.Duration(config.MaxConnectionLifeTimeMinutes)*time.Minute),
	)

	// Return connection URL.
	return url, nil
}
