package serverspace

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
	"strconv"
	"strings"
)

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordCreate,
		ReadContext:   resourceDNSRecordRead,
		DeleteContext: resourceDNSRecordDelete,
		Schema:        dnsRecordSchema,
	}
}

func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ssclient.SSClient)

	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	ip := d.Get("ip").(string)
	mailHost := d.Get("mail_host").(string)
	priority := d.Get("priority").(int)
	canonicalName := d.Get("canonical_name").(string)
	nameServerHost := d.Get("name_server_host").(string)
	text := d.Get("text").(string)
	protocol := d.Get("protocol").(string)
	service := d.Get("service").(string)
	weight := d.Get("weight").(int)
	port := d.Get("port").(int)
	target := d.Get("target").(string)
	ttl := d.Get("ttl").(string)

	dnsRecord, err := client.CreateDNSRecordAndWait(name, recordType, ip, mailHost, priority, canonicalName, nameServerHost, text, protocol, service, weight, port, target, ttl)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(dnsRecord.ID))
	resourceDNSRecordRead(ctx, d, m)
	return diags
}

func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	domain := strings.Join(strings.Split(d.Get("name").(string), ".")[1:],".")
	client := m.(*ssclient.SSClient)
	record, err := client.GetDNSRecord(domain, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", record.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", record.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ttl", record.TTL); err != nil {
		return diag.FromErr(err)
	}
	switch record.Type {
	case "A", "AAAA":
		if err := d.Set("ip", record.IP); err != nil {
			return diag.FromErr(err)
		}
	case "MX":
		if err := d.Set("mail_host", record.MailHost); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("priority", record.Priority); err != nil {
			return diag.FromErr(err)
		}
	case "SRV":
		if err := d.Set("priority", record.Priority); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("protocol", record.Protocol); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("service", record.Service); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("weight", record.Weight); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("port", record.Port); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("target", record.Target); err != nil {
			return diag.FromErr(err)
		}
	case "CNAME":
		if err := d.Set("canonical_name", record.CanonicalName); err != nil {
			return diag.FromErr(err)
		}
	case "NS":
		if err := d.Set("name_server_host", record.NameServerHost); err != nil {
			return diag.FromErr(err)
		}
	case "TXT":
		if err := d.Set("text", record.Text); err != nil {
			return diag.FromErr(err)
		}
	default:

	}

	d.SetId(strconv.Itoa(record.ID))
	return diags
}

func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*ssclient.SSClient)
	domain := strings.Join(strings.Split(d.Get("name").(string), ".")[1:],".")
	err := client.DeleteDNSRecord(domain, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
