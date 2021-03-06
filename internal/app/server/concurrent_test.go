package server_test

import (
	"testing"

	"github.com/ahamtat/itvbackend/internal/app/storage/memory"

	"github.com/ahamtat/itvbackend/internal/app/fetcher"
	"github.com/ahamtat/itvbackend/internal/app/server"
)

func TestConcurrentServer_FetchResponse(t *testing.T) {
	s := server.NewConcurrentServer(
		5,
		fetcher.NewMockFetcher(),
		memory.NewMemoryStorage())

	populateStorage(s, t)
	s.(*server.ConcurrentServer).Close()
}
