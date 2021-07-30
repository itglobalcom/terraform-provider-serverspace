package ssclient

import (
	"fmt"
)

const serverBaseURL = "servers"

type (
	VolumeData struct {
		Name   string `json:"name,omitempty"`
		SizeMB int    `json:"size_mb,omitempty"`
	}

	NetworkData struct {
		NetworkID string `json:"network_id,omitempty"`
		Bandwidth int    `json:"bandwidth_mbps,omitempty"`
	}

	ServerResponse struct {
		ID         string          `json:"id,omitempty"`
		Name       string          `json:"name,omitempty"`
		LocationID string          `json:"location_id,omitempty"`
		State      string          `json:"state,omitempty"`
		CPU        int             `json:"cpu,omitempty"`
		RAM        int             `json:"ram_mb,omitempty"`
		Volumes    []*VolumeEntity `json:"volumes,omitempty"`
		NICS       []*NICEntity    `json:"nics,omitempty"`
		SSHKeyIDS  []int           `json:"ssh_key_ids,omitempty"`
	}

	serverResponseWrap struct {
		Server *ServerResponse `json:"server,omitempty"`
	}
)

func (c *SSClient) GetServer(serverID string) (*ServerResponse, error) {
	url := fmt.Sprintf("%s/%s", serverBaseURL, serverID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &serverResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*serverResponseWrap).Server, nil
}

func (c *SSClient) CreateServer(
	name string,
	locationID string,
	imageID string,
	cpu int,
	ram int,
	volumes []*VolumeData,
	networks []*NetworkData,
	sshKeyIds []int,
) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"name":        name,
		"location_id": locationID,
		"image_id":    imageID,
		"cpu":         cpu,
		"ram_mb":      ram,
		"volumes":     volumes,
		"networks":    networks,
		"ssh_key_ids": sshKeyIds,
	}

	resp, err := makeRequest(c.client, serverBaseURL, methodPost, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateServerAndWait(
	name string,
	locationID string,
	imageID string,
	cpu int,
	ram int,
	volumes []*VolumeData,
	networks []*NetworkData,
	sshKeyIds []int,
) (*ServerResponse, error) {
	taskWrap, err := c.CreateServer(name, locationID, imageID, cpu, ram, volumes, networks, sshKeyIds)
	if err != nil {
		return nil, err
	}
	return c.waitServer(taskWrap.ID)
}

func (c *SSClient) UpdateServer(serverID string, cpu int, ram int) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"cpu":    cpu,
		"ram_mb": ram,
	}
	url := fmt.Sprintf("%s/%s", serverBaseURL, serverID)
	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) UpdateServerAndWait(serverID string, cpu int, ram int) (*ServerResponse, error) {
	taskWrap, err := c.UpdateServer(serverID, cpu, ram)
	if err != nil {
		return nil, err
	}
	return c.waitServer(taskWrap.ID)
}

func (c *SSClient) DeleteServer(serverID string) error {
	url := fmt.Sprintf("%s/%s", serverBaseURL, serverID)
	_, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{})
	return err
}

func (c *SSClient) waitServer(taskID string) (*ServerResponse, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetServer(task.ServerID)
}
