package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Headers map[string]string

type request struct {
	*Command
	location string
	headers  Headers
	method   string
	body     string
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

	return request
}

func (r *request) Name() string {
	return r.flagset.Name()
}

func (r *request) Init(args []string) {
	r.flagset.Parse(args)
}

func (r *request) parseBody() *bytes.Reader {
	// A mess
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

func buildHeaders(request *http.Request, headers Headers) {

	for header, value := range headers {
		request.Header.Add(header, value)
	}
}

func (r *request) render(response string) {

	format := `
  %s %s %s
  %s

  %s
  `
	method := r.method
	url := r.location
	scheme := "HTTP/1.1"
	var headers string

	for header, value := range r.headers {
		headers += header + ": " + value + "\n"
	}
	text := fmt.Sprintf(format, method, url, scheme, headers, response)
	fmt.Println(text)
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

	client := http.Client{}
	buildHeaders(request, parseHeaders(""))

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	r.render(string(responseBody))

	return nil
}
