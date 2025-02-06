package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Body []byte

type ResponseStatus struct {
	code   int
	status string
}

type response struct {
	headers  Headers
	body     Body
	status   ResponseStatus
	protocol string
}

func NewResponse(httpResponse *http.Response) (response response, err error) {

	response.status.code, response.status.status = httpResponse.StatusCode, httpResponse.Status
	response.body, err = io.ReadAll(httpResponse.Body)
	response.protocol = httpResponse.Proto

	if err != nil {
		return response, err
	}

	headers := Headers{
		parsedHeaders: httpResponse.Header,
		rawHeaders:    "",
	}
	response.headers = headers
	return response, nil
}

func (header *Headers) getContentType() string {
	return header.parsedHeaders["Content-Type"][0]
}

// test
func (body Body) String(contentType string) string {

	if contentType == ContentTypeJSON {

		var data interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Error parsing response body:", err)
			return ""
		}

		prettyJSON, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Error parsing response body:", err)
			return ""
		}

		return shade(ColorYellow, string(prettyJSON))
	}

	return shade(ColorYellow, string(body))
}

func (s *ResponseStatus) String() string {

	color := ColorGreen

	if s.code > 399 {
		color = ColorRed
	}

	return shade(color, s.status)
}

func (r *response) String() string {

	format := `
%s %s
%s

%s
 `
	text := fmt.Sprintf(format, r.status.String(), shade(ColorBlue, r.protocol), r.headers.String(), r.body.String(r.headers.getContentType()))
	return text
}
