package storage

import "github.com/ahamtat/itvbackend/internal/app/model"

// Storage for application requests.
type Storage interface {
	// AddFetchData saves fetch data and return ID.
	AddRequest(data *model.FetchData) (string, error)

	// AddResponse saves response from external resource by request ID.
	AddResponse(ID string, response *model.Response) error

	// GetAllRequests reads all requests from storage.
	GetAllRequests() []model.Request

	// DeleteRequest removes request from storage by ID.
	DeleteRequest(ID string) error
}
