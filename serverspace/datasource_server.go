package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

func dataSourceServerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"location": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"state": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cpu": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boot_volume_size": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boot_volume_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"volumes": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"size": {
						Type:     schema.TypeInt,
						Optional: true,
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
		"nics": {
			Type:     schema.TypeSet,
			Optional: true,
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
					},
					"bandwidth": {
						Type:     schema.TypeInt,
						Required: true,
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
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerRead,
		Schema:      dataSourceServerSchema(),
	}
}

func dataSourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	serverID := d.Get("id").(string)

	resp, err := client.GetServer(serverID)
	if err != nil {
		return diag.FromErr(err)
	}

	server := ServerToMap(resp)

	if err := d.Set("name", server["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location", server["location"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cpu", server["cpu"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ram", server["ram"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("boot_volume_size", server["boot_volume_size"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("boot_volume_id", server["boot_volume_id"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("boot_volume_id", server["boot_volume_id"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("volumes", server["volumes"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", server["tags"]); err != nil {
		return diag.FromErr(err)
	}

	d.Set("nics", server["nics"])
	d.Set("public_ip_addresses", server["public_ip_addresses"])
	d.Set("ssh_keys", server["ssh_key_ids"])

	d.SetId(serverID)

	return diags
}
