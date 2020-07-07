package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ahamtat/itvbackend/internal/app/model"
)

// DatabaseStorage data.
type DatabaseStorage struct {
	conn *Connection
	ctx  context.Context
}

// NewDatabaseStorage constructor.
func NewDatabaseStorage(ctx context.Context, conn *Connection) Storage {
	return &DatabaseStorage{
		conn: conn,
		ctx:  ctx,
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
func (s *DatabaseStorage) AddRequest(data *model.FetchData) (string, error) {
	var (
		id   string
		uuid = uuid.New().String()
	)

	// Get connection from pool
	conn, err := s.conn.Get(s.ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	err = conn.QueryRow(
		s.ctx,
		"INSERT INTO requests (uuid, method, url, fetch_headers, body) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		uuid,
		data.Method,
		data.URL,
		joinHeaders(data.Headers),
		data.Body,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

// AddResponse saves response from external resource by request ID.
func (s *DatabaseStorage) AddResponse(id string, response *model.Response) error {
	// Get connection from pool
	conn, err := s.conn.Get(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		s.ctx,
		"UPDATE requests SET status=$1, length=$2, response_headers=$3 WHERE uuid=$4",
		response.Status,
		response.Length,
		joinHeaders(response.Headers),
		id)
	return err
}

// GetAllRequests reads all requests from storage.
func (s *DatabaseStorage) GetAllRequests(paginator *model.Paginator) []model.Request {
	// TODO: Unimplemented
	return nil
}

// DeleteRequest removes request from storage by ID.
func (s *DatabaseStorage) DeleteRequest(id string) error {
	// TODO: Unimplemented
	return nil
}
