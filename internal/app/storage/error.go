package storage

import "github.com/pkg/errors"

var (
	ErrInvalidInputData = errors.New("invalid input data")
	ErrRequestNotFound  = errors.New("request not found")
)
