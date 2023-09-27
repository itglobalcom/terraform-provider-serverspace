package serverspace

import (
	"fmt"

	"github.com/itglobalcom/goss"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/utils/datalist"
)

func dataSourceImageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"os_version": {
			Type:     schema.TypeString,
			Required: true,
		},
		"architecture": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func dataSourceImages() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dataSourceImageSchema(),
		ResultAttributeName: "images",
		GetRecords:          dataSourceImagesRead,
		FlattenRecord:       flattenImage,
	}

	return datalist.NewResource(dataListConfig)
}

func dataSourceImagesRead(m interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := m.(*goss.SSClient)
	resp, err := client.GetImageList()
	if err != nil {
		return nil, err
	}

	var items []interface{}
	for _, item := range resp {
		items = append(items, *item)
	}

	return items, nil
}

func flattenImage(rawServer interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	image, ok := rawServer.(goss.ImageResponse)
	if !ok {
		return nil, fmt.Errorf("unable to convert to image")
	}

	flattenImage := map[string]interface{}{
		"type":         image.Type,
		"os_version":   image.OSVersion,
		"architecture": image.Architecture,
	}

	return flattenImage, nil
}
