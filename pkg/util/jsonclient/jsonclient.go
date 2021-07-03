package jsonclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ContentType is the default mime type for JSON
const ContentType = "application/json"

// DefaultClient is the singleton instance of jsonclient using http.DefaultClient
var DefaultClient = NewClient()

// Get fetches url using GET and unmarshals into the passed response using DefaultClient
func Get(url string, response interface{}) error {
	return DefaultClient.Get(url, response)
}

// Post sends request as JSON and unmarshals the response JSON into the supplied struct using DefaultClient
func Post(url string, request interface{}, response interface{}) error {
	return DefaultClient.Post(url, request, response)
}

// Client is a JSON wrapper around http.Client
type Client struct {
	HTTPClient          *http.Client
	AuthorizationHeader string
	Indent              string
}

func NewClient() *Client {
	return &Client{HTTPClient: http.DefaultClient}
}

// Get fetches url using GET and unmarshals into the passed response
func (c *Client) Get(url string, response interface{}) error {
	res, err := c.HTTPClient.Get(url)
	if err != nil {
		return err
	}

	return parseResponse(res, response)
}

// Post sends request as JSON and unmarshals the response JSON into the supplied struct
func (c *Client) Post(url string, request interface{}, response interface{}) error {
	var err error
	var body []byte

	body, err = json.MarshalIndent(request, "", c.Indent)
	if err != nil {
		return fmt.Errorf("error creating payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", ContentType)
	if c.AuthorizationHeader != "" {
		req.Header.Set("Authorization", c.AuthorizationHeader)
	}

	var res *http.Response
	res, err = c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending payload: %v", err)
	}

	return parseResponse(res, response)
}

func parseResponse(res *http.Response, response interface{}) error {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		err = fmt.Errorf("got HTTP %v", res.Status)
	}

	if err == nil {
		err = json.Unmarshal(body, response)
	}

	if err != nil {
		if body == nil {
			body = []byte{}
		}
		return Error{
			StatusCode: res.StatusCode,
			Body:       string(body),
			err:        err,
		}
	}

	return nil
}
