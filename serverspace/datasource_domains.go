package serverspace

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/utils/datalist"
)

func dataSourceDomains() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dataSourceDomainSchema(),
		ResultAttributeName: "domains",
		GetRecords:          dataSourceDomainsRead,
		FlattenRecord:       flattenDomain,
	}

	return datalist.NewResource(dataListConfig)
}

func dataSourceDomainsRead(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := m.(*goss.SSClient)
	resp, err := client.GetDomainList()
	if err != nil {
		return nil, err
	}

	var items []interface{}
	for _, item := range resp {
		items = append(items, *item)
	}

	return items, nil
}

func flattenDomain(rawDomain interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	domain, ok := rawDomain.(goss.DomainResponse)
	if !ok {
		return nil, fmt.Errorf("unable to convert to domain")
	}
	flattenedDomain := DomainToMap(&domain)
	return flattenedDomain, nil
}
