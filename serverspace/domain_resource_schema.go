package serverspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var domainSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}
