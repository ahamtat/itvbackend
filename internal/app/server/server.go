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
func NewServer(fetcher fetcher.Fetcher, storage storage.Storage, logger *logrus.Logger) *Server {
	// Check input data
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

	// Save request to storage
	ID, err := s.storage.AddRequest(data)
	if err != nil {
		s.logger.Errorf("makeRequest(): error saving request to storage: %s", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	// Fetch response from external resource
	resp, err := s.fetcher.Fetch(ID, data)
	if err != nil {
		s.logger.Errorf("makeRequest(): error fetching response from external resource: %s", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	// Save response to storage
	if err := s.storage.AddResponse(ID, resp); err != nil {
		s.logger.Errorf("makeRequest(): error saving response to storage: %s", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	// Return response to client
	respond(w, http.StatusOK, resp)
}

func (s *Server) deleteRequest(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		s.logger.Errorln("deleteRequest(): invalid request body")
		sendError(w, http.StatusBadRequest, nil)
		return
	}
	type request struct {
		ID string `json:"id"`
	}
	data := &request{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		s.logger.Errorf("deleteRequest(): error decoding request body: %s", err)
		sendError(w, http.StatusBadRequest, err)
		return
	}

	// Delete request from storage
	if err := s.storage.DeleteRequest(data.ID); err != nil {
		s.logger.Errorf("deleteRequest(): error saving request to storage: %s", err)
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	// Send success to client
	respond(w, http.StatusOK, nil)
}

func (s *Server) handleListAllRequests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paginator := &model.Paginator{}
		if r.Body == nil {
			paginator = nil
		} else {
			if err := json.NewDecoder(r.Body).Decode(paginator); err != nil {
				s.logger.Errorf("handleListAllRequests(): error decoding request body: %s", err)
				sendError(w, http.StatusBadRequest, err)
				return
			}
			if paginator.RequestsPerPage == 0 {
				paginator = nil
			}
		}

		// Get stored requests
		requests := s.storage.GetAllRequests(paginator)
		respond(w, http.StatusOK, requests)
	}
}
