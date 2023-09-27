package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		DeleteContext: resourceDomainDelete,
		Schema:        domainSchema,
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	migrate := d.Get("migrate_records").(bool)

	_, err := client.CreateDomainAndWait(name, migrate)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceDomainRead(ctx, d, m)

	return diags
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	domainName := d.Get("name").(string)

	domain, err := client.GetDomain(domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", domain.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_delegated", domain.IsDelegated); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domainName)

	return diags
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)
	domainName := d.Id()

	if err := client.DeleteDomain(domainName); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
