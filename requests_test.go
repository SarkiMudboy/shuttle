package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

// const headers = `{
//    "X-Frame-Options": "SAMEORIGIN",
//    "X-Content-Type-Options": "nosniff",
//    "X-Xss-Protection": [1, "mode=block]",
// }`

func TestParseHeaders(t *testing.T) {

	testCases := []struct {
		name         string
		raw          string
		resultHeader map[string][]string
		expErr       error
	}{
		{
			name: "TestNoHeaders",
			raw:  "",
			resultHeader: map[string][]string{
				"Content-Type": {"application/json"},
			},
			expErr: nil,
		},
		{
			name: "TestAddHeaderWithoutContentType",
			raw: `{
          "X-Frame-Options": "SAMEORIGIN",
          "X-Content-Type-Options": "nosniff",
          "X-Xss-Protection": ["1", "mode=block"]
      }`,
			resultHeader: map[string][]string{
				"X-Frame-Options":        {"SAMEORIGIN"},
				"X-Content-Type-Options": {"nosniff"},
				"X-Xss-Protection":       {"1", "mode=block"},
			},
			expErr: nil,
		},
		{
			name: "TestAddHeaderWithContentTypePlain",
			raw: `{
        "X-Frame-Options": "SAMEORIGIN",
        "X-Content-Type-Options": "nosniff",
        "X-Xss-Protection": ["1", "mode=block"],
"Content-Type": "text/html"
      }`,
			resultHeader: map[string][]string{
				"X-Frame-Options":        {"SAMEORIGIN"},
				"X-Content-Type-Options": {"nosniff"},
				"X-Xss-Protection":       {"1", "mode=block"},
				"Content-Type":           {ContentTypeHTML},
			},
			expErr: nil,
		},

		{
			name: "TestInvalidString",
			raw: `{
        "X-Content-Type-Options": "nosniff",
        "X-Xss-Protection": [1, "mode=block"]
      }`,
			resultHeader: make(map[string][]string),
			expErr:       ErrInvalidResourceInput, // add an error here
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			header := Headers{
				rawHeaders:    tc.raw,
				parsedHeaders: make(map[string][]string),
			}
			err := header.parse()

			if tc.expErr != nil {
				if err == nil {
					t.Error("Expected error but got none")
				}

				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error: %s, but got error: %s", tc.expErr, err)
				}

				return

			}

			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}

			if !reflect.DeepEqual(header.parsedHeaders, tc.resultHeader) {
				t.Errorf("Fail: Headers do not match. Wanted %v, got %v", tc.resultHeader, header.parsedHeaders)
			}
		})
	}
}

func TestParseRequestBody(t *testing.T) {

	testCases := []struct {
		name        string
		request     request
		contentType string
		resultBody  []byte
		expErr      error
	}{
		{
			name: "TestGetRequestWithBody",
			request: request{
				method:     "GET",
				body:       `{"houseNo":"1234","street:"New Haven"}`,
				sourceFile: "",
			},
			contentType: ContentTypeJSON,
			resultBody:  nil,
			expErr:      ErrBodyNotAllowedForSafeMethods,
		},
		{
			name: "TestGetRequestWithoutBody",
			request: request{
				method:     "GET",
				body:       "",
				sourceFile: "",
			},
			contentType: ContentTypeJSON,
			resultBody:  nil,
			expErr:      nil,
		},
		{
			name: "TestPostRequestWithBody",
			request: request{
				method:     "POST",
				body:       `{"Product_ID":333,"Product_Name":"bed"}`,
				sourceFile: "",
				Command: &Command{
					flagset: flag.NewFlagSet("", flag.ContinueOnError),
				},
			},
			contentType: ContentTypeJSON,
			resultBody:  []byte(`{"Product_ID":333,"Product_Name":"bed"}`),
			expErr:      nil,
		},
		{
			name: "TestParseRequestBodyFromFile",
			request: request{
				method:     "POST",
				sourceFile: "testdata/test_body.json",
				Command: &Command{
					flagset: flag.NewFlagSet("", flag.ContinueOnError),
				},
			},
			contentType: ContentTypeJSON,
			resultBody:  []byte(`{"Product_ID":333,"Product_Name":"bed"}`),
			expErr:      nil,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			res, err := tc.request.parseBody()

			if tc.expErr != nil {
				if err == nil {
					t.Error("Expected error got nil instead")
				}

				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error: %s, but got error: %s", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}

			if reader, ok := res.(*bytes.Reader); ok {
				result := make([]byte, reader.Len())

				reader.Seek(0, 0)
				_, err = reader.Read(result)

				if err != nil {
					t.Fatal(err)
				}

				if !bytes.Equal(tc.resultBody, bytes.TrimSpace(result)) {
					t.Errorf("Expected %s but got %s", string(tc.resultBody), string(result))
				}

			} else if tc.resultBody != nil {
				t.Fail()
			}
		})
	}
}

func TestMakeRequest(t *testing.T) {

	rd := func(method string) (r RequestData) {
		r.method = method
		r.body = `{
        "id": 1,
        "username": "jayhound101",
        "email": "jyh@gmail.com",
        "cart": {
          "no_of_items": 22,
          "cost": 233.45,
          "cleared": false
        }
      }`

		return
	}

	testCases := []struct {
		name          string
		statusCode    int
		headers       map[string][]string
		requestParams RequestData
		expErr        error
	}{
		{
			name:       "TestRequestsCanBeMadeSucessfully",
			statusCode: 200,
			headers: map[string][]string{
				"Content-Type":           {"application/json"},
				"X-Content-Type-Options": {"nosniff"},
			},
			requestParams: rd("GET"),
			expErr:        nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if slices.Contains(SafeMethods, tc.requestParams.method) {

			}
			server, err := mockServer(tc.requestParams, tc.statusCode, tc.headers)

			if err != nil {
				t.Fatalf("error creating server: %s", err)
			}

			request := NewRequest()
			request.flagset.Set("loc", server.URL)
			request.flagset.Set("method", tc.method)

			request.Init([]string{})
			fmt.Println(request.location)

			err = request.makeRequest(true)

			if tc.expErr != nil {
				if err == nil {
					t.Error("Expected error got nil instead")
				}

				if !errors.Is(tc.expErr, err) {
					t.Errorf("Expected error: %s, but got error: %s", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}

			if !headerIsSubset(request.response.headers.parsedHeaders, tc.headers) {
				t.Errorf("Headers do not match %v not %v", tc.headers, request.response.headers.parsedHeaders)
			}

			if request.response.status.code != tc.statusCode {
				t.Fatalf("Expected statusCode : %d, but got %d", request.response.status.code, tc.statusCode)
			}

			body := []byte(tc.body)

			if slices.Contains(tc.headers["Content-Type"], "application/json") {
				body, err = jsonify(tc.body)
			}

			if !bytes.Equal(bytes.TrimSpace(body), bytes.TrimSpace(request.response.body)) {
				t.Errorf("Expected %s but got %s", body, request.response.body)
			}

		})
	}
}
