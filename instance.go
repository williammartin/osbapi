package osbapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ProvisionRequest struct {
	ServiceID      string `json:"service_id"`
	PlanID         string `json:"plan_id"`
	OrganizationID string `json:"organization_guid"`
	SpaceID        string `json:"space_guid"`
}

type ProvisionResponse struct {
	DashboardURL string `json:"dashboard_url"`
	Operation    string `json:"operation"`
}

type ServiceInstance struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}

func (c *Client) Provision(instanceID string, provisionRequest *ProvisionRequest) (*ProvisionResponse, error) {
	reqBody, err := json.Marshal(provisionRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/service_instances/%s", c.brokerURL, instanceID), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-API-Version", c.apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("provision failed with status %s and body %s", resp.Status, string(body))
	}

	var provisionResponse ProvisionResponse
	if err := json.Unmarshal(body, &provisionResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body %s with error %v", string(body), err)
	}

	return &provisionResponse, nil
}

func (c *Client) GetInstance(instanceID string) (*ServiceInstance, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/service_instances/%s", c.brokerURL, instanceID), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Broker-API-Version", c.apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("get instance failed with status %s and body %s", resp.Status, string(body))
	}

	var serviceInstance ServiceInstance
	if err := json.Unmarshal(body, &serviceInstance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body %s with error %v", string(body), err)
	}

	return &serviceInstance, nil
}