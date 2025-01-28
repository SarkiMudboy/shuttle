package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Body []byte

type response struct {
	headers    Headers
	body       Body
	statusCode int
	status     string
}

func NewResponse(httpResponse *http.Response) (response response, err error) {
	// fmt.Println(httpResponse.Header)

	response.statusCode, response.status = httpResponse.StatusCode, httpResponse.Status
	response.body, err = io.ReadAll(httpResponse.Body)

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

		return string(prettyJSON)
	}

	return string(body)
}

func renderResponseStatus(code int, text string) string {

	shade := ColorGreen

	if code > 399 {
		shade = ColorRed
	}
	return fmt.Sprint(shade, text, ColorReset)
}

func (r *response) String() string {

	status := renderResponseStatus(r.statusCode, r.status)
	format := `
%s
%s

%s
  `
	text := fmt.Sprintf(format, status, r.headers.String(), r.body.String(r.headers.getContentType()))
	return text
}
