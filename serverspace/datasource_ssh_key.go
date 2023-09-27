package serverspace

import (
	"context"
	"strconv"

	"github.com/itglobalcom/goss"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSSHKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"public_key": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func dataSourceSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSSHKeyRead,
		Schema:      dataSourceSSHKeySchema(),
	}
}

func dataSourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	sshKeyID := d.Get("id").(int)

	resp, err := client.GetSSHKey(sshKeyID)
	if err != nil {
		return diag.FromErr(err)
	}

	sshKey := SSHKeyToMap(resp)

	if err := d.Set("id", sshKey["id"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", sshKey["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("public_key", sshKey["public_key"]); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(sshKeyID))

	return diags
}
