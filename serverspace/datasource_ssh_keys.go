package serverspace

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/goss"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/utils/datalist"
)

func dataSourceSSHKeys() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dataSourceSSHKeySchema(),
		ResultAttributeName: "ssh_keys",
		GetRecords:          dataSourceSSHKeysRead,
		FlattenRecord:       flattenSSHKey,
	}

	return datalist.NewResource(dataListConfig)
}

func dataSourceSSHKeysRead(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := m.(*goss.SSClient)
	resp, err := client.GetSSHKeyList()
	if err != nil {
		return nil, err
	}

	var items []interface{}
	for _, item := range resp {
		items = append(items, *item)
	}

	return items, nil
}

func flattenSSHKey(rawSSHKey interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	sshKey, ok := rawSSHKey.(goss.SSHResponse)
	if !ok {
		return nil, fmt.Errorf("unable to convert to ssh_key")
	}
	flattenedSSHKey := SSHKeyToMap(&sshKey)
	return flattenedSSHKey, nil
}
