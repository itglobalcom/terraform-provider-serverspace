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
	"boot_volume_size": {
		Type:         schema.TypeInt,
		Required:     true,
		ValidateFunc: validation.IntAtLeast(10240),
	},
	"boot_volume_id": {
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
	"public_nic": {
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"bandwidth": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 100),
				},
				"ip_address": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},

	"private_nic": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"network": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ComputedWhen: []string{"bandwidth"},
				},
				"ip_address": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			},
		},
	},
	"ssh_keys": {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		ForceNew: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
}
