package fetcher

import (
	"net/http"

	"github.com/ahamtat/itvbackend/internal/app/model"
)

type MockFetcher struct{}

// NewMockFetcher constructor.
func NewMockFetcher() Fetcher {
	return &MockFetcher{}
}

// Fetch data from mock resource.
func (f *MockFetcher) Fetch(id string, data *model.FetchData) (*model.Response, error) {
	if err := checkFetchData(data); err != nil {
		return nil, err
	}

	return &model.Response{
		ID:      id,
		Status:  http.StatusOK,
		Headers: nil,
		Length:  0,
	}, nil
}
