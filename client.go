package osbapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type ClientOpt func(c *Client)

func WithBasicAuth(username, password string) ClientOpt {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func WithAPIVersion(apiVersion string) ClientOpt {
	return func(c *Client) {
		c.apiVersion = apiVersion
	}
}

type Client struct {
	httpClient *http.Client
	brokerURL  string

	apiVersion string
	username   string
	password   string
}

func NewClient(brokerURL string, opts ...ClientOpt) *Client {
	client := &Client{
		httpClient: http.DefaultClient,
		brokerURL:  brokerURL,
		apiVersion: "2.10",
	}

	for _, o := range opts {
		o(client)
	}

	return client
}

func (c *Client) Catalog() (*Catalog, error) {
	req, err := NewRequest("GET", fmt.Sprintf("%s/v2/catalog", c.brokerURL), nil,
		WithBasicAuthHeader(c.username, c.password),
		WithAPIVersionHeader(c.apiVersion))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var catalog Catalog
	if err := json.Unmarshal(body, &catalog); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body %s with error %v", string(body), err)
	}

	return &catalog, nil
}

type RequestOpt func(r *http.Request)

func WithAPIVersionHeader(apiVersion string) RequestOpt {
	return func(req *http.Request) {
		req.Header.Set("X-Broker-API-Version", apiVersion)
	}
}

func WithContentTypeHeader() RequestOpt {
	return func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
	}
}

func WithBasicAuthHeader(username, password string) RequestOpt {
	return func(req *http.Request) {
		req.SetBasicAuth(username, password)
	}
}

func NewRequest(verb, uri string, body io.Reader, opts ...RequestOpt) (*http.Request, error) {
	req, err := http.NewRequest(verb, uri, body)
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		o(req)
	}

	return req, nil
}
