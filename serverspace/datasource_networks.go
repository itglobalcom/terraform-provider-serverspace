package serverspace

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/utils/datalist"
)

func dataSourceNetworks() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dataSourceNetworkSchema(),
		ResultAttributeName: "networks",
		GetRecords:          dataSourceNetworksRead,
		FlattenRecord:       flattenNetwork,
	}

	return datalist.NewResource(dataListConfig)
}

func dataSourceNetworksRead(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := m.(*goss.SSClient)
	resp, err := client.GetNetworkList()
	if err != nil {
		return nil, err
	}

	var items []interface{}
	for _, item := range resp {
		items = append(items, *item)
	}

	return items, nil
}

func flattenNetwork(rawNetwork interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	network, ok := rawNetwork.(goss.NetworkEntity)
	if !ok {
		return nil, fmt.Errorf("unable to convert to network")
	}
	flattenedNetwork := NetworkToMap(&network)
	return flattenedNetwork, nil
}
