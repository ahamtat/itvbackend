package model

type Paginator struct {
	// Current page number
	Page int `json:"page"`
	// Number of requests per page
	RequestsPerPage int `json:"requestsPerPage"`
}
