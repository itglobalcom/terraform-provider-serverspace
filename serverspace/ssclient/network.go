package ssclient

import (
	"fmt"
)

const networkBaseURL = "networks/isolated"

type (
	NetworkResponse struct {
		ID            string   `json:"id,omitempty"`
		Name          string   `json:"name,omitempty"`
		LocationID    string   `json:"location_id,omitempty"`
		Description   string   `json:"description,omitempty"`
		NetworkPrefix string   `json:"network_prefix,omitempty"`
		Mask          int      `json:"mask,omitempty"`
		ServerIDS     []string `json:"server_ids,omitempty"`
		State         string   `json:"state,omitempty"`
		Created       string   `json:"created,omitempty"`
	}

	networkResponseWrap struct {
		IsolatedNetwork *NetworkResponse `json:"isolated_network,omitempty"`
	}
)

func (c *SSClient) GetNetwork(networkID string) (*NetworkResponse, error) {
	url := getNetworkURL(networkID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &networkResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*networkResponseWrap).IsolatedNetwork, nil
}

func (c *SSClient) CreateNetwork(
	name string,
	locationID string,
	description string,
	networkPrefix string,
	mask int,
) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"name":           name,
		"location_id":    locationID,
		"description":    description,
		"network_prefix": networkPrefix,
		"mask":           mask,
	}
	resp, err := makeRequest(c.client, networkBaseURL, methodPost, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateNetworkAndWait(
	name string,
	locationID string,
	description string,
	networkPrefix string,
	mask int,
) (*NetworkResponse, error) {
	taskWrap, err := c.CreateNetwork(name, locationID, description, networkPrefix, mask)
	if err != nil {
		return nil, err
	}
	return c.waitNetwork(taskWrap.ID)
}

func (c *SSClient) UpdateNetwork(networkID, name, description string) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	url := getNetworkURL(networkID)
	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) UpdateNetworkAndWait(networkID, name, description string) (*NetworkResponse, error) {
	taskWrap, err := c.UpdateNetwork(networkID, name, description)
	if err != nil {
		return nil, err
	}
	return c.waitNetwork(taskWrap.ID)
}

func (c *SSClient) DeleteNetwork(networkID string) error {
	url := getNetworkURL(networkID)
	_, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{})
	return err
}

func (c *SSClient) waitNetwork(taskID string) (*NetworkResponse, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetNetwork(task.NetworkID)
}

func getNetworkURL(networkID string) string {
	return fmt.Sprintf("%s/%s", networkBaseURL, networkID)
}
