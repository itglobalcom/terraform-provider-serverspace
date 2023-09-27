package serverspace

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
)

var protectedRecords = map[string][]string{
	"NS": {
		"ns01.serverspace.com.",
		"ns02.serverspace.com.",
		"ns03.serverspace.com.",
		"ns04.serverspace.com.",
		"ns05.serverspace.com.",
		"ns06.serverspace.com.",
		"ns01.ss4test.com.",
		"ns02.ss4test.com.",
		"ns03.ss4test.com.",
		"ns04.ss4test.com.",
	},
}

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Schema:        domainRecordSchema,
	}
}

func resourceRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)
	var diags diag.Diagnostics

	recPayload, err := toRecordInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	record, err := client.CreateRecordAndWait(
		d.Get("domain").(string),
		*recPayload,
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(record.ID))
	resourceRecordRead(ctx, d, m)

	return diags
}

func resourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)
	var diags diag.Diagnostics

	domainName := d.Get("domain").(string)
	recordID := d.Get("id").(string)

	resp, err := client.GetRecord(recordID, domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	record := RecordToMap(resp)

	if err := d.Set("domain", domainName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", record["id"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", record["name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", record["type"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip", record["ip"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mail_host", record["mail_host"]); err != nil {
		return diag.FromErr(err)
	}

	if priority, ok := record["priority"]; ok && priority != 0 {
		if err := d.Set("priority", record["priority"]); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("canonical_name", record["canonical_name"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name_server_host", record["name_server_host"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("text", record["text"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("protocol", record["protocol"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("service", record["service"]); err != nil {
		return diag.FromErr(err)
	}

	if weight, ok := record["weight"]; ok && weight != 0 {
		if err := d.Set("weight", record["weight"]); err != nil {
			return diag.FromErr(err)
		}
	}

	if port, ok := record["port"]; ok && port != 0 {
		if err := d.Set("port", record["port"]); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("target", record["target"]); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ttl", record["ttl"]); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)
	var diags diag.Diagnostics

	recPayload, err := toRecordInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	record, err := client.UpdateRecordAndWait(
		d.Get("id").(string),
		d.Get("domain").(string),
		*recPayload,
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(record.ID))
	resourceRecordRead(ctx, d, m)

	return diags
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goss.SSClient)
	domainName := d.Get("domain").(string)
	recordID := d.Get("id").(string)

	if err := client.DeleteRecord(domainName, recordID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func validateProtectedRecord(recType string, recName string) error {
	for protectedecType, domains := range protectedRecords {
		if recType == protectedecType {
			for _, domain := range domains {
				if strings.HasSuffix(recName, domain) {
					return errors.New("record is protected")
				}
			}
		}
	}
	return nil
}

func toRecordInput(d *schema.ResourceData) (*goss.DomainRecord, error) {
	recName := d.Get("name").(string)
	recType := d.Get("type").(string)
	if err := validateProtectedRecord(recType, recName); err != nil {
		return nil, err
	}

	recPayload := goss.DomainRecord{
		Name: recName,
		Type: goss.RecordType(recType),
	}

	ip, ok := d.Get("ip").(string)
	if ip != "" && ok {
		recPayload.IP = &ip
	}

	mail_host, ok := d.Get("mail_host").(string)
	if mail_host != "" && ok {
		recPayload.MailHost = &mail_host
	}

	priority, ok := d.Get("priority").(int)
	if priority != -1 && ok {
		recPayload.Priority = &priority
	}

	canonical_name, ok := d.Get("canonical_name").(string)
	if canonical_name != "" && ok {
		recPayload.CanonicalName = &canonical_name
	}

	name_server_host, ok := d.Get("name_server_host").(string)
	if name_server_host != "" && ok {
		recPayload.NameServerHost = &name_server_host
	}

	text, ok := d.Get("text").(string)
	if text != "" && ok {
		recPayload.Text = &text
	}

	protocol, ok := d.Get("protocol").(string)
	if protocol != "" && ok {
		recPayload.Protocol = (*goss.ProtocolType)(&protocol)
	}

	service, ok := d.Get("service").(string)
	if service != "" && ok {
		recPayload.Service = &service
	}

	weight, ok := d.Get("weight").(int)
	if weight != -1 && ok {
		recPayload.Weight = &weight
	}

	port, ok := d.Get("port").(int)
	if port != -1 && ok {
		recPayload.Port = &port
	}

	target, ok := d.Get("target").(string)
	if target != "" && ok {
		recPayload.Target = &target
	}

	ttl, ok := d.Get("ttl").(string)
	if ttl != "" && ok {
		recPayload.TTL = &ttl
	}

	return &recPayload, nil
}
