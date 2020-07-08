package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ahamtat/itvbackend/internal/app/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // initializing postgres driver

	"github.com/google/uuid"

	"github.com/ahamtat/itvbackend/internal/app/model"
)

// Storage data.
type Storage struct {
	ctx    context.Context
	logger *logrus.Logger
	db     *sqlx.DB
}

// CreateDatabase initializes database connection pool.
func CreateDatabase(dsn string, poolSize int) (*sql.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Initialize connection pool
	db.SetMaxOpenConns(poolSize)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(0) // connections are reused forever

	return db.DB, nil
}

// NewDatabaseStorage constructor.
func NewDatabaseStorage(ctx context.Context, db *sql.DB) storage.Storage {
	return &Storage{
		ctx:    ctx,
		db:     sqlx.NewDb(db, "postgres"),
		logger: logrus.New(),
	}
}

func joinHeaders(headers map[string][]string) string {
	temp := make([]string, 0, len(headers))
	for k, v := range headers {
		temp = append(temp, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(temp, "; ")
}

// AddFetchData saves fetch data and return ID.
func (s *Storage) AddRequest(data *model.FetchData) (string, error) {
	// Create timed query context
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	var uuid = uuid.New().String()
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO requests (uuid, method, url, fetch_headers, body) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		uuid,
		data.Method,
		data.URL,
		joinHeaders(data.Headers),
		data.Body)
	if err != nil {
		s.logger.Errorf("AddRequest(): failed inserting into requests table: %s", err)
		return "", err
	}

	return uuid, nil
}

// AddResponse saves response from external resource by request ID.
func (s *Storage) AddResponse(id string, response *model.Response) error {
	// Create timed query context
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		"UPDATE requests SET status=$1, length=$2, response_headers=$3 WHERE uuid=$4",
		response.Status,
		response.Length,
		joinHeaders(response.Headers),
		id)
	if err != nil {
		s.logger.Errorf("error updating requests table: %s", err)
	}
	return err
}

// GetAllRequests reads all requests from storage.
func (s *Storage) GetAllRequests(paginator *model.Paginator) []model.Request {
	// TODO: Unimplemented
	return nil
}

// DeleteRequest removes request from storage by ID.
func (s *Storage) DeleteRequest(id string) error {
	// TODO: Unimplemented
	return nil
}
