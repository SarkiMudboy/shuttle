package main

import (
	"errors"
	"net/http"
	"strings"
)

func getMethod(method string) (string, error) {

	switch strings.ToTitle(method) {
	case "GET":
		return http.MethodGet, nil
	case "POST":
		return http.MethodPut, nil
	case "PUT":
		return http.MethodPut, nil
	case "DELETE":
		return http.MethodPut, nil
	default:
		return "", errors.New("Invalid Operation")
	}
}

func parseHeaders(rawHeaderString string) Headers {
	h := make(map[string]string)
	return h
}
