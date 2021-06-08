package ssclient

type NICResponse struct {
	ID            int    `json:"id,omitempty"`
	ServerID      string `json:"server_id,omitempty"`
	NetworkID     string `json:"network_id,omitempty"`
	MAC           string `json:"mac,omitempty"`
	IPAddress     string `json:"ip_address,omitempty"`
	Mask          int    `json:"mask,omitempty"`
	Gateway       string `json:"gateway,omitempty"`
	BandwidthMBPS int    `json:"bandwidth_mbps,omitempty"`
}
