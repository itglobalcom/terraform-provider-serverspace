package serverspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var serverSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"image": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"location": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"cpu": {
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},
	"ram": {
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validation.IntAtLeast(512),
	},
	"root_volume_size": {
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validation.IntAtLeast(10240),
	},
	"root_volume_id": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"volume": {
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"size": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(10240),
				},
			},
		},
	},
	"nic": {
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"network": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					// ConflictsWith: []string{"bandwidth"},
				},
				"network_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"bandwidth": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 100),
					// ConflictsWith: []string{"network"},
				},
			},
		},
	},
	"ssh_keys": {
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
}
