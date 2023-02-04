package serverspace

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/terraform-provider-serverspace/serverspace/ssclient"
)

func resourceSSH() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHCreate,
		ReadContext:   resourceSSHRead,
		DeleteContext: resourceSSHDelete,
		Schema:        sshSchema,
	}
}

func resourceSSHCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	publicSSH := makeNormalSSHKey(d.Get("public_key").(string))
	sshKey, err := client.CreateSSHKey(name, publicSSH)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(sshKey.ID))
	resourceSSHRead(ctx, d, m)
	return diags
}

func resourceSSHRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client := m.(*ssclient.SSClient)
	sshID, _ := strconv.Atoi(d.Id())

	sshKey, err := client.GetSSHKey(sshID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", sshKey.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("public_key", makeNormalSSHKey(sshKey.PublicKey)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(sshID))
	return diags
}

func resourceSSHDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	sshID, _ := strconv.Atoi(d.Id())

	err := client.DeleteSSHKey(sshID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func makeNormalSSHKey(sshKey string) string {
	result := sshKey
	result = strings.ReplaceAll(result, "<<~EOT", "")
	result = strings.ReplaceAll(result, "EOT", "")
	result = strings.ReplaceAll(result, "\r", "")

	lines := make([]string, 0)
	for _, line := range strings.Split(result, "\n") {
		lines = append(lines, strings.TrimSpace(line))
	}

	return strings.Join(lines, "")
}
