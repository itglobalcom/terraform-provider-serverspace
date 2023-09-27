package serverspace

import (
	"fmt"

	"github.com/itglobalcom/goss"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/utils/datalist"
)

func dataSourceLocationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func dataSourceLocations() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dataSourceLocationSchema(),
		ResultAttributeName: "locations",
		GetRecords:          dataSourceLocationsRead,
		FlattenRecord:       flattenLocation,
	}

	return datalist.NewResource(dataListConfig)
}

func dataSourceLocationsRead(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := m.(*goss.SSClient)
	resp, err := client.GetLocationList()
	if err != nil {
		return nil, err
	}

	var items []interface{}
	for _, item := range resp {
		items = append(items, *item)
	}

	return items, nil
}

func flattenLocation(rawServer interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	location, ok := rawServer.(goss.LocationEntity)
	if !ok {
		return nil, fmt.Errorf("unable to convert to location")
	}

	flattenLocation := map[string]interface{}{
		"name": location.ID,
	}

	return flattenLocation, nil
}
