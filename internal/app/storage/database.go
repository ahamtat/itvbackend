package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ahamtat/itvbackend/internal/app/model"
)

type DatabaseStorage struct {
	db *sql.DB
}

func NewDatabaseStorage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
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
	err := s.db.QueryRow(
		"INSERT INTO requests (uuid, method, url, fetch_headers, body) VALUES ($1, $2, $3, $4) RETURNING id",
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
	_, err := s.db.Exec(
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
