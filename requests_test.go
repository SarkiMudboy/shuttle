package main

import (
	//	"errors"
	"bytes"
	"errors"
	"reflect"
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
		name       string
		isFile     bool
		request    request
		resultBody []byte
		expErr     error
	}{
		{
			name:   "TestGetRequestWithBody",
			isFile: false,
			request: request{
				method: "GET",
				body:   `{"houseNo":"1234","street:"New Haven"}`,
			},
			resultBody: nil,
			expErr:     nil,
		},
		{
			name:   "TestPostRequestWithBody",
			isFile: false,
			request: request{
				method: "POST",
				body:   `{"Product_ID":333,"Product_Name":"bed"}`,
			},
			resultBody: []byte(`{"Product_ID":333,"Product_Name":"bed"}`),
			expErr:     nil,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			if !tc.isFile {
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

				result := new([]byte)
				res.Read(*result)

				// json marshal first
				if !bytes.Equal(tc.resultBody, *result) {
					t.Errorf("Expected %s but got %s", string(tc.resultBody), string(*result))
				}
			}
		})
	}
}
