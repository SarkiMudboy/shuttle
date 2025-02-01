package main

import (
	"errors"
)

const DummyEndpointTest = "https://dummyjson.com/test"

// type contentType string

const (
	ContentTypeJSON  = "application/json"
	ContentTypeHTML  = "text/html"
	ContentTypePlain = "text/plain"
)

// errors
var ErrInvalidResourceInput = errors.New("Invalid input")
