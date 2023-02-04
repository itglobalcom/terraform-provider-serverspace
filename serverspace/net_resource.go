package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/terraform-provider-serverspace/serverspace/ssclient"
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

	netwrokID := d.Id()

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err := client.DeleteNetwork(netwrokID)
		if err == nil {
			return nil
		}

		if clientErr, ok := err.(*ssclient.RequestError); ok {
			errBody := clientErr.Response.Error().(*ssclient.ErrorBodyResponse)
			if len(errBody.Errors) == 1 && errBody.Errors[0].Code == NETWORK_IS_USING_ERROR_CODE {
				return resource.RetryableError(err)
			}
		}
		return resource.NonRetryableError(err)
	})

	if err == nil {
		return diag.Diagnostics{}
	}

	return diag.FromErr(err)
}
