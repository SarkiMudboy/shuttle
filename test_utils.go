package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func mockServer(data RequestData, headers map[string][]string) (*httptest.Server, error) {

	f := func(w http.ResponseWriter, r *http.Request) {

		if slices.Contains(SafeMethods, data.method) {

			for header, value := range headers {
				values := strings.Join(value, ";")
				w.Header().Set(header, values)
			}

			if data.body != "" {
				json.NewEncoder(w).Encode(data.body)
			}
		} else {
			b, err := io.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			}

			if !bytes.Equal([]byte(data.body), b) || data.method != r.Method || !headerIsSubset(r.Header, headers) {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
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
			fmt.Println(k, v, val)
			return false
		}
	}

	return true
}
