package osbapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var INSTANCES_URL = "v2/service_instances"

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

type DeprovisionRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}

func (c *Client) Provision(instanceID string, provisionRequest *ProvisionRequest) (*ProvisionResponse, error) {
	reqBody, err := json.Marshal(provisionRequest)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s/%s/%s", c.brokerURL, INSTANCES_URL, instanceID)
	req, err := NewRequest("PUT", uri, bytes.NewBuffer(reqBody), WithCommonBrokerHeaders(c)...)
	if err != nil {
		return nil, err
	}

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

func (c *Client) Deprovision(instanceID string, deprovisionRequest *DeprovisionRequest) error {
	uri := fmt.Sprintf("%s/%s/%s", c.brokerURL, INSTANCES_URL, instanceID)
	req, err := NewRequest("DELETE", uri, nil, WithCommonBrokerHeaders(c)...)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("service_id", deprovisionRequest.ServiceID)
	q.Add("plan_id", deprovisionRequest.PlanID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("provision failed with status %s and body %s", resp.Status, string(body))
	}

	return nil
}

func (c *Client) GetInstance(instanceID string) (*ServiceInstance, error) {
	uri := fmt.Sprintf("%s/%s/%s", c.brokerURL, INSTANCES_URL, instanceID)
	req, err := NewRequest("GET", uri, nil, WithCommonBrokerHeaders(c)...)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get instance failed with status %s and body %s", resp.Status, string(body))
	}

	var serviceInstance ServiceInstance
	if err := json.Unmarshal(body, &serviceInstance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body %s with error %v", string(body), err)
	}

	return &serviceInstance, nil
}
