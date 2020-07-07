package storage

import (
	"context"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// Connection data.
type Connection struct {
	connURI  string
	pool     *pgxpool.Pool
	poolSize int
}

// NewDatabaseConnection constructor.
func NewDatabaseConnection(dsn string, poolSize int) *Connection {
	c := &Connection{
		connURI:  dsn,
		poolSize: poolSize,
	}
	return c
}

// Create and initialize connection pool to database.
func (c *Connection) Init() error {
	cfg, err := pgxpool.ParseConfig(c.connURI)
	if err != nil {
		return errors.Wrap(err, "failed to parse postgres config")
	}

	cfg.MaxConns = int32(c.poolSize)
	cfg.ConnConfig.TLSConfig = nil
	cfg.ConnConfig.PreferSimpleProtocol = true
	cfg.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on"
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		Timeout:   1 * time.Second,
		KeepAlive: 5 * time.Minute,
	}).DialContext

	// Make operation context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create connection pool
	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return err
	}
	c.pool = pool

	return nil
}

// Get connection from connection pool.
func (c *Connection) Get(ctx context.Context) (*pgxpool.Conn, error) {
	// Wrap context with timeout
	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := c.pool.Acquire(timedCtx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
