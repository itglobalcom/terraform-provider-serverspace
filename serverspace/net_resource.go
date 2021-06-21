package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Schema:        networkSchema,
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	location := d.Get("location").(string)
	netwrokProefix := d.Get("network_prefix").(string)
	mask := d.Get("mask").(int)

	network, err := client.CreateNetworkAndWait(name, location, description, netwrokProefix, mask)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(network.ID)
	resourceNetworkRead(ctx, d, m)
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)
	networkID := d.Id()

	if d.HasChanges("name", "description") {
		name := d.Get("name").(string)
		description := d.Get("description").(string)

		if _, err := client.UpdateNetwork(networkID, name, description); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*ssclient.SSClient)
	networkID := d.Id()

	network, err := client.GetNetwork(networkID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", network.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", network.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location", network.LocationID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("network_prefix", network.NetworkPrefix); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mask", network.Mask); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(networkID)
	return diags
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	netwrokID := d.Id()

	err := client.DeleteNetwork(netwrokID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
