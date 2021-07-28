package serverspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
		CustomizeDiff: customdiff.All(
			customdiff.ValidateValue("nic", func(ctx context.Context, value, meta interface{}) error {
				nics := value.(*schema.Set).List()
				for _, nic := range nics {
					mappedNIC := nic.(map[string]interface{})
					netType := ssclient.NetworkType(mappedNIC["network_type"].(string))
					if netType == ssclient.PublicSharedNetwork {
						if mappedNIC["bandwidth"].(int) == 0 {
							return fmt.Errorf("bandwidth for PublicShared interface shuold be more than 0")
						}

						if mappedNIC["network"].(string) != "" {
							return fmt.Errorf("network value for PublicShared interface shuold be \"\" (empty string)")
						}
					} else {
						if mappedNIC["bandwidth"].(int) != 0 {
							return fmt.Errorf("bandwidth for Isolated interface shuold be 0")
						}

						if mappedNIC["network"].(string) == "" {
							return fmt.Errorf("network value for Isolated interface shuold be equal some network id. Now it is  \"\" (empty string)")
						}
					}
				}
				return nil
			}),
			func(c context.Context, rd *schema.ResourceDiff, rawClient interface{}) error {
				client := rawClient.(*ssclient.SSClient)

				serverID := rd.Id()
				var err error
				// validate nics
				if serverID != "" && rd.HasChange("nic") {
					snapshots, clientErr := client.GetSnapshotList(serverID)
					if err != nil {
						return clientErr
					}
					if len(snapshots) != 0 {
						err = multierror.Append(err, fmt.Errorf("You can't change networks when have snapshots"))
					}
				}

				// validate location limits
				locationsLimit, clientErr := client.GetLocationList()
				if err != nil {
					return clientErr
				}

				location := rd.Get("location").(string)
				var locationLimit *ssclient.LocationEntity

				for _, loc := range locationsLimit {
					if loc.ID == location {
						locationLimit = loc
						break
					}
				}

				// FIXME: What should I do when location isn't found?
				if locationLimit == nil {
					return nil
				}

				// check CPU
				if rd.HasChange("cpu") {
					cpu := rd.Get("cpu").(int)

					found := false
					for _, possibleValue := range locationLimit.CPUQuantityOptions {
						if cpu == possibleValue {
							found = true
							break
						}
					}

					if !found {
						err = multierror.Append(err, fmt.Errorf("CPU value (%d) is not valid. Possible values: %v",
							cpu,
							locationLimit.CPUQuantityOptions,
						))
					}
				}

				// check RAM
				if rd.HasChange("ram") {
					ram := rd.Get("ram").(int)

					found := false
					for _, possibleValue := range locationLimit.RAMSizeOptions {
						if ram == possibleValue {
							found = true
							break
						}
					}

					if !found {
						err = multierror.Append(err, fmt.Errorf("RAM value (%d) is not valid. Possible values: %v",
							ram,
							locationLimit.RAMSizeOptions,
						))
					}
				}

				// check Boot volume
				if rd.HasChange("boot_volume_size") {
					oldRawBootSize, newRawBootSize := rd.GetChange("boot_volume_size")
					oldBootSize := oldRawBootSize.(int)
					newBootSize := newRawBootSize.(int)

					if newBootSize < locationLimit.SystemVolumeMin || newBootSize > locationLimit.VolumeMax {
						err = multierror.Append(err, fmt.Errorf(
							"boot volume size should be between %d and %d on location %s. Now it is %d",
							locationLimit.SystemVolumeMin,
							locationLimit.VolumeMax,
							location,
							newBootSize,
						))
					}
					if newBootSize < oldBootSize {
						err = multierror.Append(
							err,
							fmt.Errorf(
								"new boot volume size %d should be more than old boot size %d",
								newBootSize,
								oldBootSize,
							))
					}
					if newBootSize%1024 != 0 {
						err = multierror.Append(err, fmt.Errorf(
							"new boot volume size must be a multiple of 10Gb (1024Mb)."+
								"Now it is %d",
							newBootSize,
						))
					}
				}

				// check additional volumes
				if rd.HasChange("volume") {
					oldRawVolumes, newRawVolumes := rd.GetChange("volume")

					for i, newRawVolume := range newRawVolumes.([]interface{}) {
						newVolume := newRawVolume.(map[string]interface{})
						newVolumeSize := newVolume["size"].(int)

						if newVolumeSize < locationLimit.AdditionalVolumeMin || newVolumeSize > locationLimit.VolumeMax {
							err = multierror.Append(err, fmt.Errorf(
								"volume size should be between %d and %d in location %s. Now it is %d (index: %d)",
								locationLimit.AdditionalVolumeMin,
								locationLimit.VolumeMax,
								location,
								newVolumeSize,
								i,
							))
						}

						if newVolumeSize%1024 != 0 {
							err = multierror.Append(err, fmt.Errorf("new volume size must be a multiple of 10Gb (1024Mb). "+
								"Now it is %d",
								newVolumeSize,
							))
						}

						// However, we use list for volumes and we can check volume by position
						// https://github.com/hashicorp/terraform-plugin-sdk/issues/783
						// New size must be great then old
						oldRawVolumes := oldRawVolumes.([]interface{})
						if i <= len(oldRawVolumes)-1 {
							oldVolume := oldRawVolumes[i].(map[string]interface{})
							oldVolumeSize := oldVolume["size"].(int)

							if newVolumeSize < oldVolumeSize {
								err = multierror.Append(err, fmt.Errorf(
									"new volume size %d (index: %d) less than old volume size %d. "+
										"You can only increase volume size",
									newVolumeSize,
									i,
									oldVolumeSize,
								))
							}
						}
					}
				}

				if rd.HasChange("nic") {
					rd.SetNewComputed("public_ip_addresses")

					newNICS := rd.Get("nic").(*schema.Set).List()

					for _, newRawNIC := range newNICS {
						newNIC := newRawNIC.(map[string]interface{})
						netType := ssclient.NetworkType(newNIC["network_type"].(string))

						if netType == ssclient.PublicSharedNetwork {
							bandwidth := newNIC["bandwidth"].(int)
							if bandwidth < locationLimit.BandwidthMin || bandwidth > locationLimit.BandwidthMax {
								err = multierror.Append(err, fmt.Errorf(
									"shared network connection bandwidth should be between %d and %d in location location %s"+
										"Now it is %d",
									locationLimit.BandwidthMin,
									locationLimit.BandwidthMax,
									location,
									bandwidth,
								))
							}
						}
					}
				}

				return err
			},
		),
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*ssclient.SSClient)
	var diags diag.Diagnostics

	// ----- One value params -----

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	image := d.Get("image").(string)
	cpu := d.Get("cpu").(int)
	ram := d.Get("ram").(int)

	// ----- NICS -----

	rawNICS := d.Get("nic").(*schema.Set)
	nics := make([]*ssclient.NetworkData, rawNICS.Len())
	hasPublicSharedNIC := false

	for i, rawNIC := range rawNICS.List() {
		nic := rawNIC.(map[string]interface{})

		netType := ssclient.NetworkType(nic["network_type"].(string))
		if netType == ssclient.PublicSharedNetwork {
			nics[i] = &ssclient.NetworkData{
				Bandwidth: nic["bandwidth"].(int),
			}
			hasPublicSharedNIC = true
		} else {
			nics[i] = &ssclient.NetworkData{
				NetworkID: nic["network"].(string),
			}
		}
	}

	// ----- SSH -----

	rawSSHKeyIds := d.Get("ssh_keys").([]interface{})
	sshKeyIds := make([]int, len(rawSSHKeyIds))
	for i, v := range rawSSHKeyIds {
		sshKeyIds[i] = v.(int)
	}

	// ----- Volumes -----

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

	// ----- Root Volume -----

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

	// remove an autocreated nic if exist
	if !hasPublicSharedNIC {
		var autocreatedNIC *ssclient.NICEntity
		for _, nic := range server.NICS {
			if nic.NetworkType == ssclient.PublicSharedNetwork {
				autocreatedNIC = nic
				break
			}
		}
		if autocreatedNIC == nil {
			return diag.Errorf("can't find public shared connection but it should be")
		}
		if err := client.DeleteNIC(server.ID, autocreatedNIC.ID); err != nil {
			return diag.FromErr(err)
		}
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

		if _, err := client.UpdateServerAndWait(serverID, cpu, ram); err != nil {
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
		if _, err := client.UpdateVolumeAndWait(serverID, rootVolumeID, rootName, newRootSize); err != nil {
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

	//Aditiona volumes processing

	_, volumeChanges := d.GetChange("volume")
	rawStateVolumes := volumeChanges.([]interface{})

	stateVolumes := make([]map[string]interface{}, 0)
	for _, rawStateVolume := range rawStateVolumes {
		stateVolume := rawStateVolume.(map[string]interface{})
		if stateVolume["name"] == "boot" {
			continue
		}
		stateVolumes = append(stateVolumes, stateVolume)
	}

	sortedVolumes, err := sortVolumesByStateOrder(volumesWithoutRoot, stateVolumes)
	if err != nil {
		sortedVolumes = volumesWithoutRoot
	}

	volumes := make([]interface{}, len(sortedVolumes))
	for i, volume := range sortedVolumes {
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

	// NICS processing

	nics := make([]map[string]interface{}, 0)
	publicIPS := make([]string, 0)

	for _, nic := range server.NICS {
		var network string
		var bandwidth int
		if nic.NetworkType == ssclient.PublicSharedNetwork {
			network = ""
			bandwidth = nic.BandwidthMBPS
			publicIPS = append(publicIPS, nic.IPAddress)
		} else {
			network = nic.NetworkID
			bandwidth = 0
		}
		nics = append(nics, map[string]interface{}{
			"id":           nic.ID,
			"network_type": nic.NetworkType,
			"network":      network,
			"bandwidth":    bandwidth,
			"ip_address":   nic.IPAddress,
		})
	}
	d.Set("nic", nics)
	d.Set("public_ip_addresses", publicIPS)

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

	toCheckVolumeValues := make(map[int]map[string]interface{})
	toCreateVolumeValues := make([]map[string]interface{}, 0)

	for _, volume := range newVolumeValueIfaces.([]interface{}) {
		mappedVolume := volume.(map[string]interface{})
		volumeID := mappedVolume["id"].(int)
		if volumeID == 0 {
			toCreateVolumeValues = append(toCreateVolumeValues, mappedVolume)
		} else {
			toCheckVolumeValues[volumeID] = mappedVolume
		}
	}

	// check chenged volumes
	for volumeID, oldVolume := range oldVolumeValues {
		if volume, exist := toCheckVolumeValues[volumeID]; exist {
			size := volume["size"].(int)
			name := volume["name"].(string)
			oldSize := oldVolume["size"].(int)
			oldName := oldVolume["name"].(string)
			if size != oldSize || name != oldName {
				if _, err := client.UpdateVolumeAndWait(serverID, volumeID, name, size); err != nil {
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
	for _, newVolume := range toCreateVolumeValues {
		volumeName := newVolume["name"].(string)
		volumeSize := newVolume["size"].(int)
		if _, err := client.CreateVolumeAndWait(serverID, volumeName, volumeSize); err != nil {
			return err
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

func sortVolumesByStateOrder(
	actualVolumes []*ssclient.VolumeEntity,
	stateVolumes []map[string]interface{},
) ([]*ssclient.VolumeEntity, error) {
	newVolumesOrder := make([]*ssclient.VolumeEntity, 0)

	knownIDS := make(map[int]bool)
	for _, stateVolume := range stateVolumes {
		stateVolumeID := stateVolume["id"].(int)

		if stateVolumeID != 0 {
			knownIDS[stateVolumeID] = true
		}
	}

	for _, stateVolume := range stateVolumes {
		stateVolumeID := stateVolume["id"].(int)
		volumeMustBeNew := stateVolumeID == 0

		if volumeMustBeNew {
			foundActualValue, updatedActualVolumes, err := findNewCreatedVolume(stateVolume, actualVolumes, knownIDS)
			if err != nil {
				return nil, err
			}
			newVolumesOrder = append(newVolumesOrder, foundActualValue)
			actualVolumes = updatedActualVolumes

		} else {
			foundActualValue, updatedActualVolumes, err := findVolumeByID(stateVolumeID, actualVolumes)
			if err != nil {
				return nil, err
			}
			newVolumesOrder = append(newVolumesOrder, foundActualValue)
			actualVolumes = updatedActualVolumes
		}
	}

	return newVolumesOrder, nil
}

func findVolumeByID(
	volumeID int,
	actualVolumes []*ssclient.VolumeEntity,
) (*ssclient.VolumeEntity, []*ssclient.VolumeEntity, error) {
	for i, actualVolume := range actualVolumes {
		if volumeID == actualVolume.ID {
			copiedActualVolume := *actualVolume
			return &copiedActualVolume, removeVolumeFromSlice(actualVolumes, i), nil
		}
	}

	return nil, nil, fmt.Errorf("can't find new volume with id %d", volumeID)
}

func findNewCreatedVolume(
	targetStateVolume map[string]interface{},
	actualVolumes []*ssclient.VolumeEntity,
	knownIDS map[int]bool,
) (*ssclient.VolumeEntity, []*ssclient.VolumeEntity, error) {
	name := targetStateVolume["name"].(string)
	size := targetStateVolume["size"].(int)

	for i, actualVolume := range actualVolumes {
		isKnownID := knownIDS[actualVolume.ID]
		if !isKnownID && name == actualVolume.Name && size == actualVolume.Size {
			copiedActualVolume := *actualVolume
			return &copiedActualVolume, removeVolumeFromSlice(actualVolumes, i), nil
		}
	}

	return nil, nil, fmt.Errorf("can't find new volume with name '%s' and size %d", name, size)
}

func removeVolumeFromSlice(volumesList []*ssclient.VolumeEntity, index int) []*ssclient.VolumeEntity {
	return append(volumesList[:index], volumesList[index+1:]...)
}
