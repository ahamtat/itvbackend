package fetcher

import (
	"net/http"
	"net/url"

	"github.com/ahamtat/itvbackend/internal/app/model"
)

func checkFetchData(data *model.FetchData) error {
	if data == nil {
		return ErrInvalidInputData
	}

	// Check method type
	if m := data.Method; m != http.MethodGet && m != http.MethodPost && m != http.MethodDelete {
		return ErrWrongHTTPMethod
	}

	// Check valid URL
	if _, err := url.ParseRequestURI(data.URL); err != nil {
		return ErrInvalidInputData
	}

	return nil
}
