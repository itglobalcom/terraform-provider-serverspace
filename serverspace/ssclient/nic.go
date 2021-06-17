package ssclient

import "fmt"

type NetworkType string

const (
	PublicSharedNetwork NetworkType = "PublicShared"
	IsolatedNetwork     NetworkType = "Isolated"
)

type (
	NICEntity struct {
		ID            int         `json:"id,omitempty"`
		ServerID      string      `json:"server_id,omitempty"`
		NetworkID     string      `json:"network_id,omitempty"`
		MAC           string      `json:"mac,omitempty"`
		IPAddress     string      `json:"ip_address,omitempty"`
		Mask          int         `json:"mask,omitempty"`
		Gateway       string      `json:"gateway,omitempty"`
		BandwidthMBPS int         `json:"bandwidth_mbps,omitempty"`
		NetworkType   NetworkType `json:"network_type,omitempty"`
	}

	nicResponseWrap struct {
		NIC *NICEntity `json:"nic,omitempty"`
	}

	nicListResponseWrap struct {
		NICS []*NICEntity `json:"nics,omitempty"`
	}
)

func (c *SSClient) GetNIC(serverID string, nicID int) (*NICEntity, error) {
	url := getNICURL(serverID, nicID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &nicResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*nicResponseWrap).NIC, nil
}

func (c *SSClient) GetNICList(serverID string) ([]*NICEntity, error) {
	url := getNICSBaseURL(serverID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &nicListResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*nicListResponseWrap).NICS, nil
}

func (c *SSClient) CreateNIC(serverID, networkID string, bandwidth int) (*TaskIDWrap, error) {
	payload := make(map[string]interface{})

	if networkID != "" {
		payload["network_id"] = networkID
	} else {
		payload["bandwidth_mbps"] = bandwidth
	}

	url := getNICSBaseURL(serverID)
	resp, err := makeRequest(c.client, url, methodPost, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateNICAndWait(serverID, networkID string, bandwidth int) (*NICEntity, error) {
	taskWrap, err := c.CreateNIC(serverID, networkID, bandwidth)
	if err != nil {
		return nil, err
	}
	return c.waitNIC(serverID, taskWrap.ID)
}

func (c *SSClient) UpdatePublicNIC(serverID string, nicID, bandwidth int) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"bandwidth_mbps": bandwidth,
	}
	url := getNICURL(serverID, nicID)
	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) UpdatePublicNICAndWait(serverID string, nicID, bandwidth int) (*NICEntity, error) {
	taskWrap, err := c.UpdatePublicNIC(serverID, nicID, bandwidth)
	if err != nil {
		return nil, err
	}
	return c.waitNIC(serverID, taskWrap.ID)
}

func (c *SSClient) DeleteNIC(serverID string, nicID int) error {
	url := getNICURL(serverID, nicID)
	if _, err := makeRequest(c.client, url, methodDelete, nil, &TaskIDWrap{}); err != nil {
		return err
	}
	if _, err := c.waitServerActive(serverID); err != nil {
		return err
	}
	return nil
}

func getNICURL(serverID string, nicID int) string {
	nicBaseURL := getNICSBaseURL(serverID)
	return fmt.Sprintf("%s/%d", nicBaseURL, nicID)
}

func getNICSBaseURL(serverID string) string {
	return fmt.Sprintf("%s/%s/nics", serverBaseURL, serverID)
}

func (c *SSClient) waitNIC(serverID, taskID string) (*NICEntity, error) {
	task, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetNIC(serverID, task.NicID)
}
