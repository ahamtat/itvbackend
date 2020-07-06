package model

// FetchData from client (incoming) to external resource.
type FetchData struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers,omitempty"`
	Body    string              `json:"body,omitempty"`
}

// Response data from external resource to client (outgoing).
type Response struct {
	ID      string              `json:"id"`
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers,omitempty"`
	Length  int64               `json:"length"`
}

// Request holds incoming and outgoing data.
type Request struct {
	Fetch    *FetchData `json:"fetch"`
	Response *Response  `json:"response"`
}
