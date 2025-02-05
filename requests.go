package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

type Headers struct {
	parsedHeaders map[string][]string
	rawHeaders    string
}

type request struct {
	*Command
	response   *response
	location   string
	headers    Headers
	method     string
	body       string
	sourceFile string
}

var SafeMethods = []string{"GET", "HEAD"}

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

	request.headers.parsedHeaders = make(map[string][]string)

	return request
}

func (r *request) Name() string {
	return r.flagset.Name()
}

func (r *request) Init(args []string) {

	r.flagset.Parse(args)

	if filename := r.flagset.Arg(0); filename != "" {
		r.sourceFile = filename
	}
}

func (h *Headers) parse() (err error) {
	if h.rawHeaders != "" {

		headers := h.parsedHeaders

		headerData := make(map[string]interface{})
		err = json.Unmarshal([]byte(h.rawHeaders), &headerData)

		for key, value := range headerData {

			switch v := value.(type) {

			case []interface{}:

				var values []string
				for _, val := range v {
					if headerValue, ok := val.(string); ok {
						values = append(values, headerValue)
					} else {
						return ErrInvalidResourceInput
					}
				}

				headers[key] = values

			case string:
				headers[key] = []string{v}

			case nil:
			default:
				return ErrInvalidResourceInput
			}

		}

		if err == nil {
			h.parsedHeaders = headers
		}

	} else {
		h.parsedHeaders = defaultHeaders
	}
	return
}

// test
func (r *request) parseHeaders() (err error) {
	headers := r.headers
	err = headers.parse()

	if err == nil {
		r.headers = headers
	}
	return
}

// test
func (r *request) parseBody() (io.Reader, error) {

	//if method is not safe [] if body is not provided print warning
	// else -> if body is provided raise error else pass

	if !slices.Contains(SafeMethods, r.method) {

		if r.sourceFile != "" {
			file, err := os.Open(r.sourceFile)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			body, err := io.ReadAll(file)

			if err != nil {
				return nil, err
			}

			return bytes.NewReader(body), nil
		}

		if r.body == "" {
			log.Printf("No body provided for a %s request", r.method)
		} else {
			return bytes.NewReader([]byte(r.body)), nil
		}

	} else {
		if r.body != "" {
			return nil, ErrBodyNotAllowedForSafeMethods
		}
	}

	return nil, nil
}

func (r *request) Run() error {
	return r.makeRequest(false)
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

// test e2e
func (r *request) makeRequest(supress bool) error {

	var body io.Reader
	httpMethod, err := getMethod(r.method)

	if err != nil {
		return fmt.Errorf("An error occured: %s", err)
	}

	if r.location == "" {
		return fmt.Errorf("Enter a valid URL")
	}

	body, err = r.parseBody()
	if err != nil {
		return err
	}

	request, err := http.NewRequest(httpMethod, r.location, body) // encode the body if
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

	r.response = &res

	if !supress {
		fmt.Println(res.String())
	}
	return nil
}
