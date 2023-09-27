package serverspace

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/utils/datalist"
)

func dataSourceRecords() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dataSourceRecordSchema(),
		ResultAttributeName: "records",
		ExtraQuerySchema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		GetRecords:    dataSourceRecordsRead,
		FlattenRecord: flattenRecord,
	}

	return datalist.NewResource(dataListConfig)
}

func dataSourceRecordsRead(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := m.(*goss.SSClient)

	domain, ok := extra["domain"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `domain` key from query data")
	}

	resp, err := client.GetRecordList(domain)
	if err != nil {
		return nil, err
	}

	var items []interface{}
	for _, item := range resp {
		items = append(items, *item)
	}

	return items, nil
}

func flattenRecord(rawRecord interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	Record, ok := rawRecord.(goss.DomainRecordResponse)
	if !ok {
		return nil, fmt.Errorf("unable to convert to record")
	}

	flattenedRecord := RecordToMap(&Record)

	domain, ok := extra["domain"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `domain` key from query data")
	}
	flattenedRecord["domain"] = domain

	return flattenedRecord, nil
}
