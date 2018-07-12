package osbapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// TODO: Metadata and Schemas
type Catalog struct {
	Services []*Service `json:"services"`
}

type Service struct {
	Name                 string             `json:"name"`
	Description          string             `json:"description"`
	ID                   string             `json:"id"`
	Tags                 []string           `json:"tags"`
	Bindable             bool               `json:"bindable"`
	PlanUpdateable       bool               `json:"plan_updateable"`
	BindingsRetrievable  bool               `json:"bindings_retrievable"`
	InstancesRetrievable bool               `json:"instances_retrievable"`
	Plans                []*Plan            `json:"plans"`
	Requires             []string           `json:"requires"`
	DashboardClient      []*DashboardClient `json:"dashboard_client"`
}

type DashboardClient struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}

type Plan struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Free        bool   `json:"free"`
	Bindable    bool   `json:"bindable"`
}

func (c *Client) GetCatalog() (*Catalog, error) {
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
