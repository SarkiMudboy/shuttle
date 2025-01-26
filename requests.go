package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
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
	command.flagset.StringVar(&request.location, "loc", DummyEndpointTest, "Request's URL")

	return request
}

func (r *request) Name() string {
	return r.flagset.Name()
}

func (r *request) Init(args []string) {
	r.flagset.Parse(args)
}

func (r *request) parseBody() *bytes.Reader {
	// where our json encode will be
	return bytes.NewReader([]byte(r.body))
}

func (r *request) Run() error {
	return makeRequest(r)
}

func buildHeaders(client http.Client, headers Headers) http.Client {
	return client
}

func makeRequest(r *request) error {

	httpMethod, err := getMethod(r.method)

	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	if r.location != "" {
		return fmt.Errorf("Enter a valid URL")
	}

	request, err := http.NewRequest(httpMethod, r.location, r.parseBody()) // encode the body if

	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(responseBody))

	return nil
}
