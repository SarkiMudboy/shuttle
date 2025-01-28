package main

import (
	"errors"
	"fmt"
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

var defaultHeaders = map[string][]string{
	"Content-Type": {ContentTypeJSON},
}

func addHeadersToRequest(request *http.Request, headers Headers) {
	for header, value := range headers.parsedHeaders {
		request.Header.Set(header, strings.Join(value, ","))
	}
}

func shade(color Color, text string) string {
	return fmt.Sprint(color, text, ColorReset)
}
