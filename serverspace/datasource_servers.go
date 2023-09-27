package serverspace

import (
	"fmt"

	"github.com/itglobalcom/goss"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/utils/datalist"
)

func dataSourceServers() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dataSourceServerSchema(),
		ResultAttributeName: "servers",
		GetRecords:          dataSourceServersRead,
		FlattenRecord:       flattenServer,
	}

	return datalist.NewResource(dataListConfig)
}

func dataSourceServersRead(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := m.(*goss.SSClient)
	resp, err := client.GetServerList()
	if err != nil {
		return nil, err
	}

	var servers []interface{}
	for _, server := range resp {
		servers = append(servers, *server)
	}

	return servers, nil
}

func flattenServer(rawServer interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	server, ok := rawServer.(goss.ServerResponse)
	if !ok {
		return nil, fmt.Errorf("unable to convert to server")
	}
	flattenedServer := ServerToMap(&server)
	return flattenedServer, nil
}
