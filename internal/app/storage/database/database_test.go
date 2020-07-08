package database_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"

	"github.com/ahamtat/itvbackend/internal/app/model"

	"github.com/ahamtat/itvbackend/internal/app/storage/database"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDatabaseStorage_AddRequest(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create database storage
	s := database.NewDatabaseStorage(context.Background(), db)

	// Make database mocks
	mock.ExpectExec("INSERT INTO requests").
		WithArgs(
			sqlmock.AnyArg(),
			"GET",
			"http://google.com",
			"",
			"").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute method
	_, err = s.AddRequest(&model.FetchData{
		Method:  "GET",
		URL:     "http://google.com",
		Headers: nil,
		Body:    "",
	})
	require.Nil(t, err)

	// Make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	require.Nil(t, err)
}

func TestStorage_AddResponse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Create database storage
	s := database.NewDatabaseStorage(context.Background(), db)

	// Make database mocks
	mock.ExpectExec(
		"UPDATE requests").
		WithArgs(
			http.StatusOK,
			0,
			"",
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute method
	err = s.AddResponse(
		uuid.New().String(),
		&model.Response{
			ID:      "",
			Status:  200,
			Headers: nil,
			Length:  0,
		})
	require.Nil(t, err)

	// Make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	require.Nil(t, err)
}
