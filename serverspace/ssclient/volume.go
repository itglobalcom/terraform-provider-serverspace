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
	url := fmt.Sprintf("%s/%d", getVolumeBaseURL(serverID), volumeID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &serverResponseWrap{})
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
	url := getVolumeBaseURL(serverID)
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

func getVolumeBaseURL(serverID string) string {
	return fmt.Sprintf("%s/%s", serverBaseURL, serverID)
}

func (c *SSClient) waitVolume(serverID, taskID string) (*VolumeEntity, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetVolume(serverID, task.VolumeID)
}
