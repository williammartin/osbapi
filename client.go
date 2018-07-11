package osbapi

import (
	"encoding/json"
	"errors"
	"fmt"
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/catalog", c.brokerURL), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("X-Broker-API-Version", c.apiVersion)

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
