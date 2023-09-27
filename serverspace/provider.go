package serverspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

const userAgent string = "terraform-provider-serverspace"

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
			"serverspace_domain_record":    resourceRecord(),
		},
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"serverspace_server":         dataSourceServer(),
			"serverspace_servers":        dataSourceServers(),
			"serverspace_locations":      dataSourceLocations(),
			"serverspace_images":         dataSourceImages(),
			"serverspace_network":        dataSourceNetwork(),
			"serverspace_networks":       dataSourceNetworks(),
			"serverspace_ssh_key":        dataSourceSSHKey(),
			"serverspace_ssh_keys":       dataSourceSSHKeys(),
			"serverspace_domain":         dataSourceDomain(),
			"serverspace_domains":        dataSourceDomains(),
			"serverspace_domain_record":  dataSourceRecord(),
			"serverspace_domain_records": dataSourceRecords(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	key := d.Get("key").(string)
	host := d.Get("host").(string)
	customUserAgent := userAgent

	var diags diag.Diagnostics

	c, err := goss.NewClient(key, host, &customUserAgent)
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
