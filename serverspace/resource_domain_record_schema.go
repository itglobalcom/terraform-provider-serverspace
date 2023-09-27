package serverspace

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var domainRecordSchema = map[string]*schema.Schema{
	"domain": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringMatch(regexp.MustCompile("^.+\\.$"), "domain name must ends with dot"),
	},
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringMatch(regexp.MustCompile("^.+\\.$"), "record name must ends with dot"),
	},
	"type": {
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"A",
			"AAAA",
			"MX",
			"CNAME",
			"MX",
			"NS",
			"TXT",
			"SRV",
		}, false),
	},
	"ip": {
		Type: schema.TypeString,
		Optional: true,
	},
	"mail_host": {
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringMatch(regexp.MustCompile("^.+\\.$"), "record mail host must ends with dot"),
	},
	"priority": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      -1,
		ValidateFunc: validation.IntBetween(0, 65535),
	},
	"canonical_name": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"name_server_host": {
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringMatch(regexp.MustCompile("^.+\\.$"), "record name server host must ends with dot"),
	},
	"text": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"protocol": {
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"tcp",
			"udp",
			"tls",
		}, false),
	},
	"service": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"weight": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      -1,
		ValidateFunc: validation.IntBetween(0, 65535),
	},
	"port": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      -1,
		ValidateFunc: validation.IntBetween(0, 65535),
	},
	"target": {
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringMatch(regexp.MustCompile("^.+\\.$"), "record target must ends with dot"),
	},
	"ttl": {
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice([]string{
			"1s",
			"5s",
			"30s",
			"1m",
			"5m",
			"10m",
			"15m",
			"30m",
			"1h",
			"2h",
			"6h",
			"12h",
			"1d",
		}, false),
	},
}
