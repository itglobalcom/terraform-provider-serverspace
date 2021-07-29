package serverspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
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
		ValidateFunc: validation.IntAtLeast(1),
	},
	"boot_volume_size": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"boot_volume_id": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"volume": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 4,
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
	"public_ip_addresses": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"nic": {
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 5,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"network": {
					Type:     schema.TypeString,
					Required: true,
				},
				"network_type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(ssclient.PublicSharedNetwork),
						string(ssclient.IsolatedNetwork),
					}, false),
				},
				"bandwidth": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 1024),
				},
				"ip_address": {
					Type:     schema.TypeString,
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
