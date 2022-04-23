package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVERSPACE_HOST", nil),
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVERSPACE_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"serverspace_server":           resourceServer(),
			"serverspace_isolated_network": resourceNetwork(),
			"serverspace_ssh":              resourceSSH(),
			"serverspace_domain":           resourceDomain(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	key := d.Get("key").(string)
	host := d.Get("host").(string)

	var diags diag.Diagnostics

	c, err := ssclient.NewClient(key, host)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create serverspace client",
			Detail:   err.Error(), // "Unable to authenticate user for authenticated serverspace client",
		})

		return nil, diags
	}

	return c, diags

}
