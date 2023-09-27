package serverspace

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var domainSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringMatch(regexp.MustCompile("^.+\\.$"), "domain name must ends with dot"),
	},
	"migrate_records": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		ForceNew: true,
	},
	"is_delegated": {
		Type:     schema.TypeBool,
		Computed: true,
	},
}
