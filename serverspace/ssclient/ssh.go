package ssclient

import (
	"fmt"
)

const sshBaseURL = "ssh-keys"

type (
	SSHResponse struct {
		ID        int    `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		PublicKey string `json:"public_key,omitempty"`
	}

	sshResponseWrap struct {
		SSHKey *SSHResponse `json:"ssh_key,omitempty"`
	}
)

func (c *SSClient) GetSSHKey(sshID int) (*SSHResponse, error) {
	url := fmt.Sprintf("%s/%d", sshBaseURL, sshID)
	resp, err := makeRequest(c.client, url, methodGet, nil, &sshResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*sshResponseWrap).SSHKey, nil
}

func (c *SSClient) CreateSSHKey(
	name string,
	publicKey string,
) (*SSHResponse, error) {
	payload := map[string]interface{}{
		"name":       name,
		"public_key": publicKey,
	}

	resp, err := makeRequest(c.client, sshBaseURL, methodPost, payload, &SSHResponse{})
	if err != nil {
		return nil, err
	}
	return resp.(*SSHResponse), nil
}

func (c *SSClient) DeleteSSHKey(sshID int) error {
	url := fmt.Sprintf("%s/%d", sshBaseURL, sshID)
	_, err := makeRequest(c.client, url, methodDelete, nil, nil)
	return err
}
