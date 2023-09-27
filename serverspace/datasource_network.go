package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

func dataSourceNetworkSchema() map[string]*schema.Schema {
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
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"network_prefix": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mask": {
			Type:     schema.TypeInt,
			Optional: true,
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

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkRead,
		Schema:      dataSourceNetworkSchema(),
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	networkID := d.Get("id").(string)

	resp, err := client.GetNetwork(networkID)
	if err != nil {
		return diag.FromErr(err)
	}

	network := NetworkToMap(resp)

	if err := d.Set("name", network["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", network["description"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location", network["location"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("network_prefix", network["network_prefix"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mask", network["mask"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", network["tags"]); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(networkID)

	return diags
}
