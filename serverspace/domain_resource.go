package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
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
	client := m.(*ssclient.SSClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	domain, err := client.CreateDomainAndWait(name, false)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain.Name)

	resourceDomainRead(ctx, d, m)
	return diags
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*ssclient.SSClient)

	domain, err := client.GetDomain(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", domain.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)

	err := client.DeleteDomain(name)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
