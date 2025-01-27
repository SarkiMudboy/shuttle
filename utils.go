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
		return http.MethodPost, nil
	case "PUT":
		return http.MethodPut, nil
	case "DELETE":
		return http.MethodDelete, nil
	default:
		return "", errors.New("Invalid Operation")
	}
}

func parseHeaders(rawHeaderString string) Headers {
	h := make(map[string]string)

	if rawHeaderString == "" { // default headers
		h["Content-Type"] = "application/json"
	}
	return h
}
