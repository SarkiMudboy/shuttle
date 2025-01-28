package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Headers struct {
	parsedHeaders map[string][]string
	rawHeaders    string
}

type request struct {
	*Command
	response *response
	location string
	headers  Headers
	method   string
	body     string
}

func (h *Headers) String() (headers string) {

	for header, value := range h.parsedHeaders {
		headerValue := strings.Join(value, ",")
		headers += (header + ": " + headerValue + "\n")
	}
	return
}

func NewRequest() *request {

	command := &Command{
		flagset: flag.NewFlagSet("call", flag.ContinueOnError),
	}
	request := &request{Command: command}

	// flags
	command.flagset.StringVar(&request.location, "loc", DummyEndpointTest, "Define request's URL")
	command.flagset.StringVar(&request.method, "method", "GET", "Define request's HTTP Method")
	command.flagset.StringVar(&request.body, "data", "", "Define the raw data for the request's body")
	command.flagset.StringVar(&request.headers.rawHeaders, "headers", "", "Add headers for the request")

	return request
}

func (r *request) Name() string {
	return r.flagset.Name()
}

func (r *request) Init(args []string) {
	r.flagset.Parse(args)
}

func (r *request) parseHeaders() (err error) {

	if r.headers.rawHeaders != "" {

		headers := r.headers.parsedHeaders
		err = json.Unmarshal([]byte(r.headers.rawHeaders), &headers)

		if err == nil {
			r.headers.parsedHeaders = headers
		}

	} else {
		r.headers.parsedHeaders = defaultHeaders
	}
	return
}

func (r *request) parseBody() io.Reader {
	// A mess
	// update: not a mess

	if (r.method == "POST" || r.method == "PUT" || r.method == "PATCH") && r.body == "" {

		if filename := r.flagset.Arg(0); filename != "" {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()

			body, err := io.ReadAll(file)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			return bytes.NewReader(body)
		}
	} else if r.body != "" {
		return bytes.NewReader([]byte(r.body))
	} else {
		fmt.Printf("No body provided for a %s request", r.method)
	}

	return nil
}

func (r *request) Run() error {
	return makeRequest(r)
}

func (r *request) String() string {

	format := `
%s %s %s
%s
%s
  `
	var body Body

	method := r.method
	url := r.location
	scheme := "HTTP/1.1"
	body = []byte(r.body)

	text := fmt.Sprintf(format, method, url, scheme, r.headers.String(), body.String(r.headers.getContentType()))
	return text
}

func makeRequest(r *request) error {

	httpMethod, err := getMethod(r.method)

	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	if r.location == "" {
		return fmt.Errorf("Enter a valid URL")
	}

	request, err := http.NewRequest(httpMethod, r.location, r.parseBody()) // encode the body if
	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	err = r.parseHeaders()
	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	addHeadersToRequest(request, r.headers)

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	defer response.Body.Close()

	res, err := NewResponse(response)
	if err != nil {
		return err
	}

	fmt.Println(res.String())
	return nil
}
