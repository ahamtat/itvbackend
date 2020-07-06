package storage

import (
	"sync"

	"github.com/ahamtat/itvbackend/internal/app/model"
	"github.com/google/uuid"
)

// MemoryStorage makes memory implementation of Storage interface.
type MemoryStorage struct {
	mx      sync.Mutex
	storage map[string]*model.Request
}

// NewMemoryStorage constructor.
func NewMemoryStorage() Storage {
	return &MemoryStorage{
		mx:      sync.Mutex{},
		storage: make(map[string]*model.Request),
	}
}

// AddFetchData saves fetch data and return ID.
func (s *MemoryStorage) AddRequest(data *model.FetchData) (string, error) {
	if data == nil {
		return "", ErrInvalidInputData
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	// Create new request in memory
	ID := uuid.New().String()
	s.storage[ID] = &model.Request{
		Fetch:    data,
		Response: nil,
	}
	return ID, nil
}

// AddResponse saves response from external resource by request ID.
func (s *MemoryStorage) AddResponse(id string, response *model.Response) error {
	// Check input data
	if response == nil {
		return ErrInvalidInputData
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	// Save response in memory
	req, ok := s.storage[id]
	if !ok {
		return ErrRequestNotFound
	}
	req.Response = response
	return nil
}

// GetAllRequests reads all requests from storage.
func (s *MemoryStorage) GetAllRequests(paginator *model.Paginator) []model.Request {
	capacity := len(s.storage)
	if paginator != nil {
		capacity = paginator.RequestsPerPage
	}
	result := make([]model.Request, 0, capacity)

	s.mx.Lock()
	defer s.mx.Unlock()

	// Copy requests for reliability
	index := -1
	for _, value := range s.storage {
		index++

		// Skip request from undesirable page
		if paginator != nil && index/paginator.RequestsPerPage != paginator.Page {
			continue
		}

		result = append(result, model.Request{
			Fetch:    value.Fetch,
			Response: value.Response,
		})
	}
	return result
}

// DeleteRequest removes request from storage by ID.
func (s *MemoryStorage) DeleteRequest(id string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	// Remove response from memory
	_, ok := s.storage[id]
	if !ok {
		return ErrRequestNotFound
	}
	delete(s.storage, id)
	return nil
}
