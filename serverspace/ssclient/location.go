package ssclient

type (
	LocationEntity struct {
		ID                     string `json:"id,omitempty"`
		SystemVolumeMin        int    `json:"system_volume_min,omitempty"`
		AdditionalVolumeMin    int    `json:"additional_volume_min,omitempty"`
		VolumeMax              int    `json:"volume_max,omitempty"`
		WindowsSystemVolumeMin int    `json:"windows_system_volume_min,omitempty"`
		BandwidthMin           int    `json:"bandwidth_min,omitempty"`
		BandwidthMax           int    `json:"bandwidth_max,omitempty"`
		CPUQuantityOptions     []int  `json:"cpu_quantity_options,omitempty"`
		RAMSizeOptions         []int  `json:"ram_size_options,omitempty"`
	}

	LocationEntityWrap struct {
		Location *LocationEntity `json:"location,omitempty"`
	}

	locationListResponseWrap struct {
		Locations []*LocationEntity `json:"locations,omitempty"`
	}
)

func (c *SSClient) GetLocationList() ([]*LocationEntity, error) {
	url := getLocationBaseURL()
	resp, err := makeRequest(c.client, url, methodGet, nil, &locationListResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*locationListResponseWrap).Locations, nil
}

func getLocationBaseURL() string {
	return "locations"
}
