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
	"nic": {
		Type:     schema.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"network": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ComputedWhen: []string{"bandwidth"},
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
					Optional:     true,
					Computed:     true,
					ComputedWhen: []string{"network"},
					ValidateFunc: validation.IntBetween(0, 100),
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
