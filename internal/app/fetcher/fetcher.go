package fetcher

import "github.com/ahamtat/itvbackend/internal/app/model"

// Fetcher interface for external resource.
type Fetcher interface {
	// Fetch data from external resource.
	Fetch(ID string, data *model.FetchData) (*model.Response, error)
}
