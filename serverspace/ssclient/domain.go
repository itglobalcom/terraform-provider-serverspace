package ssclient

import (
	"fmt"
)

const domainBaseURL = "domains"

type (
	DomainResponse struct {
		Name        string `json:"name,omitempty"`
		IsDelegated bool   `json:"is_delegated,omitempty"`
	}

	domainResponseWrap struct {
		Domain *DomainResponse `json:"domain,omitempty"`
	}
)

func (c *SSClient) GetDomain(domainName string) (*DomainResponse, error) {
	url := fmt.Sprintf("%s/%s", domainBaseURL, domainName)
	resp, err := makeRequest(c.client, url, methodGet, nil, &domainResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*domainResponseWrap).Domain, nil
}

func (c *SSClient) CreateDomain(name string, migrate bool) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"name":            name,
		"migrate_records": migrate,
	}

	resp, err := makeRequest(c.client, domainBaseURL, methodPost, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateDomainAndWait(name string, migrate bool) (*DomainResponse, error) {
	taskWrap, err := c.CreateDomain(name, migrate)
	if err != nil {
		return nil, err
	}
	return c.waitDomain(taskWrap.ID)
}

func (c *SSClient) DeleteDomain(name string) error {
	url := fmt.Sprintf("%s/%s", domainBaseURL, name)
	_, err := makeRequest(c.client, url, methodDelete, nil, nil)
	return err
}

func (c *SSClient) waitDomain(taskID string) (*DomainResponse, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetDomain(task.DomainID)
}
