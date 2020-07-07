package fetcher

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ahamtat/itvbackend/internal/app/model"
)

// HTTPFetcher for external resource implements Fetcher interface.
type HTTPFetcher struct {
	client *http.Client
}

// NewHTTPFetcher constructor.
func NewHTTPFetcher(timeout time.Duration) Fetcher {
	return &HTTPFetcher{client: &http.Client{
		Timeout: timeout,
	}}
}

// Fetch data from external resource.
func (f *HTTPFetcher) Fetch(id string, data *model.FetchData) (*model.Response, error) {
	if err := checkFetchData(data); err != nil {
		return nil, err
	}

	// Create HTTP request to external resource
	var body io.Reader
	if len(data.Body) > 0 {
		body = bytes.NewReader([]byte(data.Body))
	}
	req, err := http.NewRequest(data.Method, data.URL, body)
	if err != nil {
		return nil, ErrCreatingHTTPRequest
	}

	// Proxying HTTP headers to request
	for key, value := range data.Headers {
		req.Header.Add(key, strings.Join(value, " "))
	}

	// Make request to external resource
	resp, err := f.client.Do(req)
	if err != nil || resp == nil {
		// Process error from external resource
		statusCode := http.StatusInternalServerError
		if resp != nil {
			statusCode = resp.StatusCode
		}
		return &model.Response{
			ID:      id,
			Status:  statusCode,
			Headers: nil,
			Length:  0,
		}, nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	length := resp.ContentLength
	if length < 0 {
		length = 0
	}

	// Get info from valid response
	return &model.Response{
		ID:      id,
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Length:  length,
	}, nil
}
