package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahamtat/itvbackend/internal/app/storage/memory"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/ahamtat/itvbackend/internal/app/server"

	"github.com/ahamtat/itvbackend/internal/app/fetcher"
	"github.com/ahamtat/itvbackend/internal/app/model"
	"github.com/stretchr/testify/require"
)

var fetchData = []model.FetchData{
	{
		Method:  "GET",
		URL:     "http://google.com",
		Headers: nil,
		Body:    "",
	},
	{
		Method:  "GET",
		URL:     "http://google.com",
		Headers: nil,
		Body:    "",
	},
	{
		Method:  "GET",
		URL:     "http://google.com",
		Headers: nil,
		Body:    "",
	},
}

func populateStorage(s http.Handler, t *testing.T) {
	// Populate storage with responses
	for _, d := range fetchData {
		body, err := json.Marshal(d)
		require.Nil(t, err)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/v1/requests/request", bytes.NewReader(body))
		require.Nil(t, err)
		s.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestServer_FetchResponse(t *testing.T) {
	s := server.NewServer(
		fetcher.NewMockFetcher(),
		memory.NewMemoryStorage(),
		logrus.New())

	populateStorage(s, t)
}

func readAndDecodeRequests(s http.Handler, expected int, paginator *model.Paginator, t *testing.T) []model.Request {
	var body io.Reader
	if paginator != nil && paginator.RequestsPerPage > 0 {
		buff, err := json.Marshal(paginator)
		require.Nil(t, err)
		body = bytes.NewReader(buff)
	}
	// Read responses from storage
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/v1/requests/list", body)
	require.Nil(t, err)
	s.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, rec.Body)

	// Decode and test responses
	var result []model.Request
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	require.Nil(t, err)
	require.Equal(t, expected, len(result))

	return result
}

func TestServer_ListResponse(t *testing.T) {
	s := server.NewServer(
		fetcher.NewMockFetcher(),
		memory.NewMemoryStorage(),
		logrus.New())

	// Read All data
	populateStorage(s, t)
	_ = readAndDecodeRequests(s, len(fetchData), nil, t)

	// Read with paginator
	populateStorage(s, t)
	_ = readAndDecodeRequests(s, 2, &model.Paginator{
		Page:            0,
		RequestsPerPage: 2,
	}, t)
}

func TestServer_DeleteResponse(t *testing.T) {
	s := server.NewServer(
		fetcher.NewMockFetcher(),
		memory.NewMemoryStorage(),
		logrus.New())

	populateStorage(s, t)
	requests := readAndDecodeRequests(s, len(fetchData), nil, t)

	// Delete responses from storage
	for _, req := range requests {
		// Create body with ID
		type requestBody struct {
			ID string `json:"id"`
		}
		reqBody := &requestBody{ID: req.Response.ID}
		assert.NotEmpty(t, reqBody.ID)
		body, err := json.Marshal(reqBody)
		require.Nil(t, err)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/v1/requests/request", bytes.NewReader(body))
		require.Nil(t, err)
		s.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}

	// Check if storage is empty
	emptyRequest := readAndDecodeRequests(s, 0, nil, t)
	require.Empty(t, emptyRequest)
}
