package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

func dataSourceDomainSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"is_delegated": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema:      dataSourceDomainSchema(),
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	domainName := d.Get("name").(string)

	resp, err := client.GetDomain(domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	domain := DomainToMap(resp)

	if err := d.Set("name", domain["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_delegated", domain["is_delegated"]); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domainName)

	return diags
}

func DomainToMap(domain *goss.DomainResponse) map[string]interface{} {
	return map[string]interface{}{
		"name":         domain.Name,
		"is_delegated": domain.IsDelegated,
	}
}
