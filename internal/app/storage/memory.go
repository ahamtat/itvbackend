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
func (s *MemoryStorage) AddResponse(ID string, response *model.Response) error {
	// Check input data
	if response == nil {
		return ErrInvalidInputData
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	// Save response in memory
	req, ok := s.storage[ID]
	if !ok {
		return ErrRequestNotFound
	}
	req.Response = response
	return nil
}

// GetAllRequests reads all requests from storage.
func (s *MemoryStorage) GetAllRequests() []model.Request {
	result := make([]model.Request, 0, len(s.storage))

	s.mx.Lock()
	defer s.mx.Unlock()

	// Copy requests for reliability
	for _, v := range s.storage {
		result = append(result, model.Request{
			Fetch:    v.Fetch,
			Response: v.Response,
		})
	}
	return result
}

// DeleteRequest removes request from storage by ID.
func (s *MemoryStorage) DeleteRequest(ID string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	// Remove response from memory
	_, ok := s.storage[ID]
	if !ok {
		return ErrRequestNotFound
	}
	delete(s.storage, ID)
	return nil
}
