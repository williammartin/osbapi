package osbapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type BindingRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}

type UnbindingRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}

type ServiceBinding struct {
	Credentials interface{} `json:"credentials"`
}

func (c *Client) Bind(instanceID, bindingID string, bindingRequest *BindingRequest) (*ServiceBinding, error) {
	reqBody, err := json.Marshal(bindingRequest)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s/%s/%s/service_bindings/%s", c.brokerURL, INSTANCES_URL, instanceID, bindingID)
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
		return nil, fmt.Errorf("binding failed with status %s and body %s", resp.Status, string(body))
	}

	var serviceBinding ServiceBinding
	if err := json.Unmarshal(body, &serviceBinding); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body %s with error %v", string(body), err)
	}

	return &serviceBinding, nil
}

func (c *Client) Unbind(instanceID, bindingID string, unbindingRequest *UnbindingRequest) error {
	uri := fmt.Sprintf("%s/%s/%s/service_bindings/%s", c.brokerURL, INSTANCES_URL, instanceID, bindingID)
	req, err := NewRequest("DELETE", uri, nil, WithCommonBrokerHeaders(c)...)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("service_id", unbindingRequest.ServiceID)
	q.Add("plan_id", unbindingRequest.PlanID)
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
		return fmt.Errorf("unbinding failed with status %s and body %s", resp.Status, string(body))
	}

	return nil
}

func (c *Client) GetBinding(instanceID, bindingID string) (*ServiceBinding, error) {
	uri := fmt.Sprintf("%s/%s/%s/service_bindings/%s", c.brokerURL, INSTANCES_URL, instanceID, bindingID)
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

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("get binding failed with status %s and body %s", resp.Status, string(body))
	}

	var serviceBinding ServiceBinding
	if err := json.Unmarshal(body, &serviceBinding); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body %s with error %v", string(body), err)
	}

	return &serviceBinding, nil
}
