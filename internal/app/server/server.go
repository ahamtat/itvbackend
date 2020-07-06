package server

import (
	"encoding/json"
	"net/http"

	"github.com/ahamtat/itvbackend/internal/app/fetcher"
	"github.com/ahamtat/itvbackend/internal/app/model"
	"github.com/ahamtat/itvbackend/internal/app/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server holds data for application logic
type Server struct {
	router  *mux.Router
	logger  *logrus.Logger
	fetcher fetcher.Fetcher
	storage storage.Storage
}

// NewServer constructor.
func NewServer(fetcher fetcher.Fetcher, storage storage.Storage) *Server {
	// Check input data
	logger := logrus.New()
	if fetcher == nil || storage == nil {
		logger.Fatalf("NewServer(): invalid input data")
	}

	s := &Server{
		router:  mux.NewRouter(),
		logger:  logger,
		fetcher: fetcher,
		storage: storage,
	}

	s.configureRouter()
	return s
}

// ServeHTTP implementation for external handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) configureRouter() {
	requests := s.router.PathPrefix("/v1/requests").Subrouter()
	requests.HandleFunc("/request", s.handleRequest()).Methods("POST", "DELETE")
	requests.HandleFunc("/list", s.handleListAllRequests()).Methods("GET")
}

func (s *Server) handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.makeRequest(w, r)
		case http.MethodDelete:
			s.deleteRequest(w, r)
		}
	}
}

func (s *Server) makeRequest(w http.ResponseWriter, r *http.Request) {
	data := &model.FetchData{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.Errorf("makeRequest(): error decoding request body: %s", err)
		s.error(w, http.StatusBadRequest, err)
		return
	}

	// Save request to storage
	ID, err := s.storage.AddRequest(data)
	if err != nil {
		s.logger.Errorf("makeRequest(): error saving request to storage: %s", err)
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	// Fetch response from external resource
	resp, err := s.fetcher.Fetch(ID, data)
	if err != nil {
		s.logger.Errorf("makeRequest(): error fetching response from external resource: %s", err)
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	// Save response to storage
	if err := s.storage.AddResponse(ID, resp); err != nil {
		s.logger.Errorf("makeRequest(): error saving response to storage: %s", err)
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	// Return response to client
	s.respond(w, http.StatusOK, resp)
}

func (s *Server) deleteRequest(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ID string `json:"id"`
	}
	data := &request{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.Errorf("deleteRequest(): error decoding request body: %s", err)
		s.error(w, http.StatusBadRequest, err)
		return
	}

	// Delete request from storage
	if err := s.storage.DeleteRequest(data.ID); err != nil {
		s.logger.Errorf("deleteRequest(): error saving request to storage: %s", err)
		s.error(w, http.StatusInternalServerError, err)
		return
	}

	// Send success to client
	s.respond(w, http.StatusOK, nil)
}

func (s *Server) handleListAllRequests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get stored requests
		requests := s.storage.GetAllRequests()
		s.respond(w, http.StatusOK, requests)
	}
}

func (s *Server) error(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			s.logger.Errorf("failed encoding JSON: %v", err)
		}
	}
}
