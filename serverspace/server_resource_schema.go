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
	"volume": {
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"size": {
					Type:     schema.TypeInt,
					Required: true,
				},
			},
		},
	},
	"nics": {
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
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
