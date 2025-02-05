package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
)

type RequestData struct {
	body   string
	method string
}

func mockServer(data RequestData, status int, headers map[string][]string) (*httptest.Server, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(status)

		if slices.Contains(SafeMethods, data.method) {

			for header, value := range headers {
				values := strings.Join(value, ";")
				w.Header().Set(header, values)
			}

			if data.body != "" {
				json.NewEncoder(w).Encode(data.body)
			}
		} // else block that check request props against the provided vars and returns a 400 and writes error to response Writer

	}

	return httptest.NewServer(http.HandlerFunc(f)), nil
}

func jsonify(body string) ([]byte, error) {

	j, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func headerIsSubset(responseHeaders, testHeaders map[string][]string) bool {

	if len(testHeaders) > len(responseHeaders) {
		return false
	}

	for k, v := range testHeaders {
		if val, ok := responseHeaders[k]; !ok || !reflect.DeepEqual(val, v) {
			return false
		}
	}

	return true
}
