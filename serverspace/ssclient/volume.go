package ssclient

import (
	"fmt"
)

type (
	VolumeEntity struct {
		ID      int    `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Size    int    `json:"size_mb,omitempty"`
		Created string `json:"created,omitempty"`
	}

	volumeResponseWrap struct {
		Volume *VolumeEntity `json:"volume,omitempty"`
	}
)

func (c *SSClient) GetVolume(serverID string, volumeID int) (*VolumeEntity, error) {
	volumeBaseURL := getVolumesBaseURL(serverID)
	url := fmt.Sprintf("%s/%d", volumeBaseURL, volumeID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &volumeResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*volumeResponseWrap).Volume, nil
}

func (c *SSClient) CreateVolume(serverID, name string, size int) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"server_id": serverID,
		"name":      name,
		"size_mb":   size,
	}
	url := getVolumesBaseURL(serverID)
	resp, err := makeRequest(c.client, url, methodPost, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateVolumeAndWait(
	serverID string,
	name string,
	size int,
) (*VolumeEntity, error) {
	taskWrap, err := c.CreateVolume(serverID, name, size)
	if err != nil {
		return nil, err
	}
	return c.waitVolume(serverID, taskWrap.ID)
}

func (c *SSClient) UpdateVolume(
	serverID string,
	volumeID int,
	name string,
	size int,
) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"name":    name,
		"size_mb": size,
	}
	url := getVolumeURL(serverID, volumeID)
	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) UpdateVolumeAndWait(
	serverID string,
	volumeID int,
	name string,
	size int,
) (*VolumeEntity, error) {
	taskWrap, err := c.UpdateVolume(serverID, volumeID, name, size)
	if err != nil {
		return nil, err
	}
	return c.waitVolume(serverID, taskWrap.ID)
}

func (c *SSClient) DeleteVolume(serverID string, volumeID int) error {
	url := getVolumeURL(serverID, volumeID)
	if _, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{}); err != nil {
		return err
	}
	if _, err := c.waitServerActive(serverID); err != nil {
		return err
	}
	return nil
}

func getVolumeURL(serverID string, volumeID int) string {
	volumesBaseURL := getVolumesBaseURL(serverID)
	return fmt.Sprintf("%s/%d", volumesBaseURL, volumeID)
}

func getVolumesBaseURL(serverID string) string {
	return fmt.Sprintf("%s/%s/volumes", serverBaseURL, serverID)
}

func (c *SSClient) waitVolume(serverID, taskID string) (*VolumeEntity, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetVolume(serverID, task.VolumeID)
}
