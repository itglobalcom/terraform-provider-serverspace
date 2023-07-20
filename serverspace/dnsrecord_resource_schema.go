package serverspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var dnsRecordSchema = map[string]*schema.Schema{
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
	"type": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"A",
			"AAAA",
			"MX",
			"CNAME",
			"NS",
			"TXT",
			"SRV",
		}, false),
	},
	"ip": {
		Type:         schema.TypeString,
		Optional: true,
		ForceNew: true,
		ValidateFunc: validation.IsIPAddress,
	},
	"mail_host": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"priority": {
		Type:         schema.TypeInt,
		Optional:     true,
		ForceNew: true,
		ValidateFunc: validation.IntBetween(0, 65535),
	},
	"canonical_name": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"name_server_host": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"text": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"protocol": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"TCP",
			"UDP",
			"TLS",
		}, false),
	},
	"service": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"weight": {
		Type:         schema.TypeInt,
		Optional:     true,
		ForceNew: true,
		ValidateFunc: validation.IntBetween(0, 65535),
	},
	"port": {
		Type:         schema.TypeInt,
		Optional:     true,
		ForceNew: true,
		ValidateFunc: validation.IntBetween(1, 65535),
	},
	"target": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"ttl": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
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
