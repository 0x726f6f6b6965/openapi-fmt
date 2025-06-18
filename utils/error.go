package utils

import "errors"

var (
	ErrOpenAPINotFound     = errors.New("OpenAPI document not found")
	ErrOpenAPIPathNotFound = errors.New("OpenAPI path not found")
)
