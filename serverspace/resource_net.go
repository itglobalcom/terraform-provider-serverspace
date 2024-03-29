package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

const NETWORK_IS_USING_ERROR_CODE = -19511

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
	client := m.(*goss.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	location := d.Get("location").(string)
	networkPrefix := d.Get("network_prefix").(string)
	mask := d.Get("mask").(int)

	network, err := client.CreateNetworkAndWait(name, location, description, networkPrefix, mask)
	if err != nil {
		return diag.FromErr(err)
	}

	// tag network
	if err := client.TagNetwork(network.ID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(network.ID)
	resourceNetworkRead(ctx, d, m)
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)
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

	client := m.(*goss.SSClient)
	networkID := d.Id()

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

	d.SetId(networkID)
	return diags
}

func NetworkToMap(network *goss.NetworkEntity) map[string]interface{} {
	networkMap := map[string]interface{}{
		"id":             network.ID,
		"name":           network.Name,
		"location":       network.LocationID,
		"description":    network.Description,
		"network_prefix": network.NetworkPrefix,
		"mask":           network.Mask,
		"tags":           network.Tags,
	}
	return networkMap
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)

	networkID := d.Id()

	err := resource.RetryContext(
		context.Background(),
		d.Timeout(schema.TimeoutDelete),
		func() *resource.RetryError {
			err := client.DeleteNetwork(networkID)
			if err == nil {
				return nil
			}

			if clientErr, ok := err.(*goss.RequestError); ok {
				errBody := clientErr.Response.Error().(*goss.ErrorBodyResponse)
				if len(errBody.Errors) == 1 && errBody.Errors[0].Code == NETWORK_IS_USING_ERROR_CODE {
					return resource.RetryableError(err)
				}
			}
			return resource.NonRetryableError(err)
		},
	)

	if err == nil {
		return diag.Diagnostics{}
	}

	return diag.FromErr(err)
}
