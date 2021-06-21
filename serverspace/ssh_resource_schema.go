package serverspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sshSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
		ForceNew: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"public_key": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}
