package serverspace

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

func dataSourceRecordSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain": {
			Type:     schema.TypeString,
			Required: true,
		},
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mail_host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"priority": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"canonical_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name_server_host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"text": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"service": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"weight": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"target": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ttl": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecordRead,
		Schema:      dataSourceRecordSchema(),
	}
}

func dataSourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	domainName := d.Get("domain").(string)
	recordID := d.Get("id").(string)

	resp, err := client.GetRecord(recordID, domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	record := RecordToMap(resp)

	if err := d.Set("domain", domainName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", record["id"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", record["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", record["type"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip", record["ip"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mail_host", record["mail_host"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("priority", record["priority"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("canonical_name", record["canonical_name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name_server_host", record["name_server_host"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("text", record["text"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("protocol", record["protocol"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("service", record["service"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("weight", record["weight"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("port", record["port"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("target", record["target"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ttl", record["ttl"]); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(recordID)

	return diags
}

func RecordToMap(record *goss.DomainRecordResponse) map[string]interface{} {
	recordMap := map[string]interface{}{
		"id":   strconv.Itoa(record.ID),
		"name": record.Name,
		"type": record.Type,
	}

	if record.IP != nil {
		recordMap["ip"] = *record.IP
	}

	if record.MailHost != nil {
		recordMap["mail_host"] = *record.MailHost
	}

	if record.Priority != nil {
		recordMap["priority"] = *record.Priority
	}

	if record.CanonicalName != nil {
		recordMap["canonical_name"] = *record.CanonicalName
	}

	if record.NameServerHost != nil {
		recordMap["name_server_host"] = *record.NameServerHost
	}

	if record.Text != nil {
		recordMap["text"] = *record.Text
	}

	if record.Protocol != nil {
		recordMap["protocol"] = *record.Protocol
	}

	if record.Service != nil {
		recordMap["service"] = *record.Service
	}

	if record.Weight != nil {
		recordMap["weight"] = *record.Weight
	}

	if record.Port != nil {
		recordMap["port"] = *record.Port
	}

	if record.Target != nil {
		recordMap["target"] = *record.Target
	}

	if record.TTL != nil {
		recordMap["ttl"] = *record.TTL
	}

	return recordMap
}
