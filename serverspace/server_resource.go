package serverspace

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Schema:        serverSchema,
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	image := d.Get("image").(string)
	cpu := d.Get("cpu").(int)
	ram := d.Get("ram").(int)

	netBandwidths := d.Get("nics").([]interface{})
	nics := make([]*ssclient.Network, len(netBandwidths), len(netBandwidths))
	for i, v := range netBandwidths {
		nics[i] = &ssclient.Network{
			Bandwidth: v.(int),
		}
	}

	rawSSHKeyIds := d.Get("ssh_keys").([]interface{})
	sshKeyIds := make([]int, len(rawSSHKeyIds), len(rawSSHKeyIds))
	for i, v := range rawSSHKeyIds {
		sshKeyIds[i] = v.(int)
	}

	rawVolumes := d.Get("volume").([]interface{})
	volumes := make([]*ssclient.VolumeData, len(rawVolumes), len(rawVolumes))
	for i, v := range rawVolumes {
		rawVolume := v.(map[string]interface{})
		volume := &ssclient.VolumeData{
			Name:   rawVolume["name"].(string),
			SizeMB: rawVolume["size"].(int),
		}
		volumes[i] = volume
	}

	server, err := client.CreateServerAndWait(name, location, image, cpu, ram, volumes, nics, sshKeyIds)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(server.ID)

	resourceServerRead(ctx, d, m)

	return diags
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	serverID := d.Get("id").(string)

	if d.HasChanges("cpu", "ram") {
		cpu := d.Get("cpu").(int)
		ram := d.Get("ram").(int)

		if _, err := client.UpdateServer(serverID, cpu, ram); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("volume") {
		oldVolumeValueIfaces, newVolumeValueIfaces := d.GetChange("volume")

		oldVolumeValues := make(map[int]map[string]interface{}, len(oldVolumeValueIfaces.([]interface{})))
		for _, volume := range oldVolumeValueIfaces.([]map[string]interface{}) {
			oldVolumeValues[volume["id"].(int)] = volume
		}

		newVolumeValues := make(map[int]map[string]interface{}, len(newVolumeValueIfaces.([]interface{})))
		for _, volume := range newVolumeValueIfaces.([]map[string]interface{}) {
			volumeID := volume["id"].(int)
			newVolumeValues[volumeID] = volume
		}

		for oldVolumeID, oldVolume := range oldVolumeValues {
			if newVolume, exist := newVolumeValues[oldVolumeID]; exist {
				if newVolume["size"].(int) != oldVolume["size"].(int) {
					_ = newVolume["size"].(int)
				}
			}
		}

	}

	return resourceServerRead(ctx, d, m)
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	serverID := d.Get("id").(string)

	server, err := client.GetServer(serverID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", server.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location", server.LocationID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cpu", server.CPU); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ram", server.RAM); err != nil {
		return diag.FromErr(err)
	}

	volumes := make([]interface{}, len(server.Volumes), len(server.Volumes))
	for i, volume := range server.Volumes {
		volumeMap := map[string]interface{}{
			"id":   volume.ID,
			"name": volume.Name,
			"size": volume.Size,
		}
		volumes[i] = volumeMap
	}
	if err := d.Set("volume", volumes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serverID)
	log.Default().Printf("999999999999 %+v", d)

	return diags
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderID := d.Id()

	err := client.DeleteServer(orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
