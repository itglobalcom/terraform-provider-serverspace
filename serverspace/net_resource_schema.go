package serverspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var networkSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"location": {
		Type:     schema.TypeString,
		Required: true,
	},
	"description": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"mask": {
		Type:     schema.TypeInt,
		Required: true,
		ForceNew: true,
	},
}
