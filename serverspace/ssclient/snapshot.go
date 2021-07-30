package ssclient

import "fmt"

type (
	SnapshotEntity struct {
		ID       int    `json:"id,omitempty"`
		ServerID string `json:"server_id,omitempty"`
		Name     string `json:"name,omitempty"`
		SizeMB   int    `json:"size_mb,omitempty"`
		Created  string `json:"created,omitempty"`
	}

	SnapshotEntityWrap struct {
		Snapshot *SnapshotEntity `json:"snapshot,omitempty"`
	}

	snapshotListResponseWrap struct {
		Snapshots []*SnapshotEntity `json:"snapshots,omitempty"`
	}
)

func (c *SSClient) GetSnapshotList(serverID string) ([]*SnapshotEntity, error) {
	url := getSnapshotBaseURL(serverID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &snapshotListResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*snapshotListResponseWrap).Snapshots, nil
}

func getSnapshotURL(serverID string, snapshotID int) string {
	snapshotBaseURL := getSnapshotBaseURL(serverID)
	return fmt.Sprintf("%s/%d", snapshotBaseURL, snapshotID)
}

func getSnapshotBaseURL(serverID string) string {
	return fmt.Sprintf("%s/%s/snapshots", serverBaseURL, serverID)
}
