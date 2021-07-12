package serverspace

import (
	"context"

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
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	image := d.Get("image").(string)
	cpu := d.Get("cpu").(int)
	ram := d.Get("ram").(int)

	rawPublicNICS := d.Get("public_nic").(*schema.Set).List()
	rawPrivateNICS := d.Get("private_nic").(*schema.Set).List()
	nics := make([]*ssclient.NetworkData, len(rawPublicNICS)+len(rawPrivateNICS))

	for i, rawNic := range rawPublicNICS {
		nic := rawNic.(map[string]interface{})
		nics[i] = &ssclient.NetworkData{
			Bandwidth: nic["bandwidth"].(int),
		}
	}

	for i, rawNic := range rawPrivateNICS {
		nic := rawNic.(map[string]interface{})
		nics[i] = &ssclient.NetworkData{
			NetworkID: nic["network"].(string),
		}
	}

	rawSSHKeyIds := d.Get("ssh_keys").([]interface{})
	sshKeyIds := make([]int, len(rawSSHKeyIds))
	for i, v := range rawSSHKeyIds {
		sshKeyIds[i] = v.(int)
	}

	// ----- Set Volumes -----
	rawVolumes := d.Get("volume").([]interface{})
	volumes := make([]*ssclient.VolumeData, len(rawVolumes))

	for i, v := range rawVolumes {
		rawVolume := v.(map[string]interface{})
		volume := &ssclient.VolumeData{
			Name:   rawVolume["name"].(string),
			SizeMB: rawVolume["size"].(int),
		}
		volumes[i] = volume
	}

	// ----- Set Root Volume -----
	rootVolumeSize := d.Get("boot_volume_size").(int)
	rootVolume := &ssclient.VolumeData{
		Name:   "boot",
		SizeMB: rootVolumeSize,
	}
	volumes = append(volumes, rootVolume)

	// ----- Perform server creating -----
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
		if err := updateVolumes(d, client, serverID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("boot_volume_size") {
		rootVolumeID := d.Get("boot_volume_id").(int)
		newRootSize := d.Get("boot_volume_size").(int)
		rootName := "boot"
		if _, err := client.UpdateVolume(serverID, rootVolumeID, rootName, newRootSize); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("nic") {
		if err := updateNICS(d, client, serverID); err != nil {
			return diag.FromErr(err)
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

	volumesWithoutRoot, rootVolume := splitRootFromVolumes(server.Volumes)
	if err := d.Set("boot_volume_size", rootVolume.Size); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("boot_volume_id", rootVolume.ID); err != nil {
		return diag.FromErr(err)
	}

	volumes := make([]interface{}, len(volumesWithoutRoot))
	for i, volume := range volumesWithoutRoot {
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

	publicNICS := make([]map[string]interface{}, 0)
	privateNICS := make([]map[string]interface{}, 0)
	for _, nic := range server.NICS {
		if nic.NetworkType == ssclient.PublicSharedNetwork {
			publicNICS = append(publicNICS, map[string]interface{}{
				"id":        nic.ID,
				"bandwidth": nic.BandwidthMBPS,
			})
		} else {
			privateNICS = append(privateNICS, map[string]interface{}{
				"network": nic.NetworkID,
			})
		}
	}
	d.Set("public_nic", publicNICS)
	d.Set("private_nic", privateNICS)

	d.Set("ssh_keys", server.SSHKeyIDS)

	d.SetId(serverID)

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

	d.SetId("")

	return diags
}

func updateVolumes(d *schema.ResourceData, client *ssclient.SSClient, serverID string) error {
	oldVolumeValueIfaces, newVolumeValueIfaces := d.GetChange("volume")

	oldVolumeValues := make(map[int]map[string]interface{}, len(oldVolumeValueIfaces.([]interface{})))
	for _, volume := range oldVolumeValueIfaces.([]interface{}) {
		mappedVolume := volume.(map[string]interface{})
		volumeID := mappedVolume["id"].(int)
		oldVolumeValues[volumeID] = mappedVolume
	}

	newVolumeValues := make(map[int]map[string]interface{}, len(newVolumeValueIfaces.([]interface{})))
	for _, volume := range newVolumeValueIfaces.([]interface{}) {
		mappedVolume := volume.(map[string]interface{})
		volumeID := mappedVolume["id"].(int)
		newVolumeValues[volumeID] = mappedVolume
	}

	// ----- VOLUMES -----
	// check chenged volumes
	for volumeID, oldVolume := range oldVolumeValues {
		if newVolume, exist := newVolumeValues[volumeID]; exist {
			newSize := newVolume["size"].(int)
			newName := newVolume["name"].(string)
			oldSize := oldVolume["size"].(int)
			oldName := oldVolume["name"].(string)
			if newSize != oldSize || newName != oldName {
				if _, err := client.UpdateVolume(serverID, volumeID, newName, newSize); err != nil {
					return err
				}
			}
		} else {
			// if volume was removed
			if err := client.DeleteVolume(serverID, volumeID); err != nil {
				return err
			}
		}
	}

	// try to find new volumes
	for volumeID, newVolume := range newVolumeValues {
		if _, exist := oldVolumeValues[volumeID]; !exist {
			volumeName := newVolume["name"].(string)
			volumeSize := newVolume["size"].(int)
			if _, err := client.CreateVolumeAndWait(serverID, volumeName, volumeSize); err != nil {
				return err
			}
		}
	}
	return nil
}

func splitRootFromVolumes(volumes []*ssclient.VolumeEntity) ([]*ssclient.VolumeEntity, *ssclient.VolumeEntity) {
	volumesLen := len(volumes)
	volumesWithoutRoot := make([]*ssclient.VolumeEntity, 0, volumesLen)
	var rootVolume *ssclient.VolumeEntity
	for _, volume := range volumes {
		if volume.Name == "boot" {
			rootVolume = volume
			continue
		}
		volumesWithoutRoot = append(volumesWithoutRoot, volume)
	}

	return volumesWithoutRoot, rootVolume
}
