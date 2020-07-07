package server

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ahamtat/itvbackend/internal/app/fetcher"
	"github.com/ahamtat/itvbackend/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/ahamtat/itvbackend/internal/app/model"
)

// ConcurrentServer data
type ConcurrentServer struct {
	router  *mux.Router
	logger  *logrus.Logger
	fetcher fetcher.Fetcher
	storage storage.Storage

	poolSize int
	taskCh   chan *model.FetchData
	Wg       sync.WaitGroup
}

// NewConcurrentServer constructor.
func NewConcurrentServer(poolSize int, fetcher fetcher.Fetcher, storage storage.Storage, logger *logrus.Logger) *ConcurrentServer {
	s := &ConcurrentServer{
		router:   mux.NewRouter(),
		fetcher:  fetcher,
		storage:  storage,
		logger:   logger,
		poolSize: poolSize,
		taskCh:   make(chan *model.FetchData, poolSize),
	}
	s.configureRouter()

	// Create workers
	for i := 0; i < poolSize; i++ {
		go s.worker()
	}
	return s
}

func (s *ConcurrentServer) worker() {
	for data := range s.taskCh {
		// Save request to storage
		id, err := s.storage.AddRequest(data)
		if err != nil {
			s.logger.Errorf("worker(): error saving request to storage: %s", err)
			return
		}

		// Fetch response from external resource
		resp, err := s.fetcher.Fetch(id, data)
		if err != nil {
			s.logger.Errorf("worker(): error fetching response from external resource: %s", err)
			return
		}

		// Save response to storage
		if err := s.storage.AddResponse(id, resp); err != nil {
			s.logger.Errorf("worker(): error saving response to storage: %s", err)
			return
		}

		//fmt.Println("worker ended")
		s.Wg.Done()
	}
}

// ServeHTTP implementation for external handler.
func (s *ConcurrentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *ConcurrentServer) configureRouter() {
	requests := s.router.PathPrefix("/v1/requests").Subrouter()
	requests.HandleFunc("/request", s.handleRequest()).Methods("POST")
}

func (s *ConcurrentServer) handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.makeRequest(w, r)
		default:
			sendError(w, http.StatusBadRequest, nil)
		}
	}
}

func (s *ConcurrentServer) makeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		s.logger.Errorln("makeRequest(): invalid request body")
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	data := &model.FetchData{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.Errorf("makeRequest(): error decoding request body: %s", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	// Send data to task channel
	s.Wg.Add(1)
	s.taskCh <- data
}
