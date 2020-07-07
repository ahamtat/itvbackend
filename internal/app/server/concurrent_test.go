package server_test

import (
	"testing"

	"github.com/ahamtat/itvbackend/internal/app/fetcher"
	"github.com/ahamtat/itvbackend/internal/app/server"
	"github.com/ahamtat/itvbackend/internal/app/storage"
	"github.com/sirupsen/logrus"
)

func TestConcurrentServer_FetchResponse(t *testing.T) {
	s := server.NewConcurrentServer(
		5,
		fetcher.NewMockFetcher(),
		storage.NewMemoryStorage(),
		logrus.New())

	populateStorage(s, t)
	s.Wg.Wait()
}
