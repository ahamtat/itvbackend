package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ahamtat/itvbackend/internal/app/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_AddRequest(t *testing.T) {
	testCases := []struct {
		name             string
		fetch            *model.FetchData
		notEmptyExpected bool
		errExpected      error
	}{
		{
			name: "Empty data",
			fetch: &model.FetchData{
				Method:  "",
				URL:     "",
				Headers: nil,
				Body:    "",
			},
			notEmptyExpected: true,
			errExpected:      nil,
		},
		{
			name:             "Nil data",
			fetch:            nil,
			notEmptyExpected: false,
			errExpected:      ErrInvalidInputData,
		},
	}

	s := NewMemoryStorage()
	require.NotNil(t, s)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ID, err := s.AddRequest(tc.fetch)
			require.Equal(t, tc.notEmptyExpected, len(ID) > 0)
			require.Equal(t, tc.errExpected, err)
		})
	}
}

func TestMemoryStorage_AddResponse(t *testing.T) {
	s := NewMemoryStorage()
	require.NotNil(t, s)

	// Populate storage with data
	generatedID := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		ID, err := s.AddRequest(&model.FetchData{
			Method:  "GET",
			URL:     "http://google.com",
			Headers: nil,
			Body:    "",
		})
		require.NotEmpty(t, ID)
		require.Nil(t, err)
		// Save generated IDs
		generatedID = append(generatedID, ID)
	}
	require.Equal(t, len(generatedID), 10)

	// Add response to existing request
	err := s.AddResponse(generatedID[5], &model.Response{
		ID:      generatedID[5],
		Status:  http.StatusOK,
		Headers: nil,
		Length:  0,
	})
	require.Nil(t, err)

	// Add response to non-existing request
	fakeID := uuid.New().String()
	err = s.AddResponse(fakeID, &model.Response{
		ID:      fakeID,
		Status:  http.StatusInternalServerError,
		Headers: nil,
		Length:  0,
	})
	require.NotNil(t, err)
}

func TestMemoryStorage_GetAllRequests(t *testing.T) {
	s := NewMemoryStorage()
	require.NotNil(t, s)

	// Populate storage with data
	generatedID := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		ID, err := s.AddRequest(&model.FetchData{
			Method:  "GET",
			URL:     "http://google.com",
			Headers: nil,
			Body:    "",
		})
		require.NotEmpty(t, ID)
		require.Nil(t, err)
		// Save generated IDs
		generatedID = append(generatedID, ID)
	}
	require.Equal(t, len(generatedID), 10)

	// Get request list
	requests := s.GetAllRequests()
	require.Equal(t, len(requests), 10)
	for _, req := range requests {
		assert.Equal(t, &model.Request{
			Fetch: &model.FetchData{
				Method:  "GET",
				URL:     "http://google.com",
				Headers: nil,
				Body:    "",
			},
			Response: nil,
		}, &req)
	}
}

func TestMemoryStorage_DeleteRequest(t *testing.T) {
	s := NewMemoryStorage()
	require.NotNil(t, s)

	// Populate storage with data
	generatedID := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		ID, err := s.AddRequest(&model.FetchData{
			Method:  "GET",
			URL:     "http://google.com",
			Headers: nil,
			Body:    "",
		})
		require.NotEmpty(t, ID)
		require.Nil(t, err)
		// Save generated IDs
		generatedID = append(generatedID, ID)
	}
	require.Equal(t, len(generatedID), 10)

	// Delete some requests by existing ID
	require.Nil(t, s.DeleteRequest(generatedID[0]))
	require.Nil(t, s.DeleteRequest(generatedID[5]))
	require.Nil(t, s.DeleteRequest(generatedID[9]))

	// Delete request by invalid ID
	require.Equal(t, ErrRequestNotFound, s.DeleteRequest(uuid.New().String()))
}
