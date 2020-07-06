package fetcher

import "github.com/pkg/errors"

var (
	ErrInvalidInputData    = errors.New("invalid input data")
	ErrCreatingHTTPRequest = errors.New("error creating HTTP request")
	ErrWrongHTTPMethod     = errors.New("wrong HTTP method")
)
