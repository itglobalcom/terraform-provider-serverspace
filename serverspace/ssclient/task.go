package ssclient

import (
	"fmt"
	"log"
	"time"
)

const defaultTaskCompletionDuration = 5 * time.Minute

type (
	TaskIDWrap struct {
		ID string `json:"task_id,omitempty"`
	}

	TaskResponse struct {
		ID          string `json:"id,omitempty"`
		Created     string `json:"created,omitempty"`
		Completed   string `json:"completed,omitempty"`
		IsCompleted string `json:"is_completed,omitempty"`

		ServerID   string `json:"server_id,omitempty"`
		LocationID string `json:"location_id,omitempty"`
		NetworkID  string `json:"network_id,omitempty"`
		VolumeID   int    `json:"volume_id,omitempty"`
		NicID      int    `json:"nic_id,omitempty"`
		SnapshotID int    `json:"snapshot_id,omitempty"`
		DomainID   string `json:"domain_id,omitempty"`
		RecordID   int    `json:"record_id,omitempty"`
	}

	taskResponseWrap struct {
		Task *TaskResponse `json:"task,omitempty"`
	}
)

func (c *SSClient) GetTask(taskID string) (*TaskResponse, error) {
	url := fmt.Sprintf("tasks/%s", taskID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &taskResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*taskResponseWrap).Task, nil
}

func (c *SSClient) waitTaskCompletion(taskID string) (*TaskResponse, error) {
	const duration = defaultTaskCompletionDuration
	begin := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var (
		task *TaskResponse
		err  error
	)

	for range ticker.C {
		task, err = c.GetTask(taskID)
		if err != nil {
			return nil, err
		}
		if task.IsCompleted == "Completed" {
			return task, nil
		} else if task.IsCompleted == "Failed" {
			return nil, fmt.Errorf("Task '%s' failed", task.ID)

		} else {
			log.Default().Printf("[TRACE] Task isn't completed: %#v", task)
		}
		if time.Now().Sub(begin) > duration {
			return nil, fmt.Errorf("Task wasn't complete for %f secs", duration.Seconds())
		}
	}

	return task, err
}

func (c *SSClient) waitServerActive(serverID string) (*ServerResponse, error) {
	const duration = defaultTaskCompletionDuration
	begin := time.Now()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var (
		server *ServerResponse
		err    error
	)

	for range ticker.C {
		server, err = c.GetServer(serverID)
		if err != nil {
			return nil, err
		}
		if server.State == "Active" {
			return server, nil
		} else {
			log.Default().Printf("[TRACE] Server isn't active: %#v", server)
		}
		if time.Now().Sub(begin) > duration {
			return nil, fmt.Errorf("Server wasn't active for %f secs", duration.Seconds())
		}
	}

	return server, err
}
