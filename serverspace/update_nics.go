package serverspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.itglobal.com/b2c/terraform-provider-serverspace/serverspace/ssclient"
)

func updateNICS(d *schema.ResourceData, client *ssclient.SSClient, serverID string) error {
	oldPublicNICS := make([]map[string]interface{}, 0)
	oldPrivateNICS := make([]map[string]interface{}, 0)
	newPublicNICS := make([]map[string]interface{}, 0)
	newPrivateNICS := make([]map[string]interface{}, 0)

	oldNICSValueIfaces, newNICSValueIfaces := d.GetChange("nic")

	for _, nic := range oldNICSValueIfaces.([]interface{}) {
		mappedNIC := nic.(map[string]interface{})
		netType := ssclient.NetworkType(mappedNIC["network_type"].(string))
		if netType == ssclient.PublicSharedNetwork {
			oldPublicNICS = append(oldPublicNICS, mappedNIC)
		} else {
			oldPrivateNICS = append(oldPrivateNICS, mappedNIC)
		}
	}

	for _, nic := range newNICSValueIfaces.([]interface{}) {
		mappedNIC := nic.(map[string]interface{})
		netType := ssclient.NetworkType(mappedNIC["network_type"].(string))
		if netType == ssclient.PublicSharedNetwork {
			newPublicNICS = append(newPublicNICS, mappedNIC)
		} else {
			newPrivateNICS = append(newPrivateNICS, mappedNIC)
		}
	}

	if err := updatePrivateNICS(client, serverID, oldPrivateNICS, newPrivateNICS); err != nil {
		return err
	}

	return updatePublicNICS(client, serverID, oldPublicNICS, newPublicNICS)
}

func updatePublicNICS(
	client *ssclient.SSClient,
	serverID string,
	oldNICS, newNICS []map[string]interface{},
) error {
	for _, oldNIC := range oldNICS {
		nicID := oldNIC["id"].(int)
		oldBandwidth := oldNIC["bandwidth"].(int)

		if newNIC := findNICByID(newNICS, nicID); newNIC == nil {
			if err := client.DeleteNIC(serverID, nicID); err != nil {
				return err
			}
		} else {
			newBandwidth := newNIC["bandwidth"].(int)
			if oldBandwidth != newBandwidth {
				if _, err := client.UpdatePublicNICAndWait(serverID, nicID, newBandwidth); err != nil {
					return err
				}
			}
		}
	}
	for _, newNIC := range newNICS {
		nicID := newNIC["id"].(int)
		if findNICByID(oldNICS, nicID) == nil {
			newNICBandwidth := newNIC["bandwidth"].(int)
			if _, err := client.CreateNICAndWait(serverID, "", newNICBandwidth); err != nil {
				return err
			}
		}
	}
	return nil
}

func updatePrivateNICS(
	client *ssclient.SSClient,
	serverID string,
	oldNICS, newNICS []map[string]interface{},
) error {
	currentNICEntities, err := client.GetNICList(serverID)
	if err != nil {
		return err
	}

	for _, oldNIC := range oldNICS {
		networkID := oldNIC["network"].(string)
		if findNICByNetwork(newNICS, networkID) == nil {
			nicEntity := findPrivateNICEntityByNetwork(currentNICEntities, networkID)
			if nicEntity == nil {
				continue
			}
			if err := client.DeleteNIC(serverID, nicEntity.ID); err != nil {
				return err
			}
		}
	}
	for _, newNIC := range newNICS {
		newNetworkID := newNIC["network"].(string)
		if findNICByNetwork(oldNICS, newNetworkID) == nil {
			if _, err := client.CreateNICAndWait(serverID, newNetworkID, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

func findPrivateNICEntityByNetwork(nics []*ssclient.NICEntity, networkID string) *ssclient.NICEntity {
	for _, nic := range nics {
		if nic.NetworkType == ssclient.IsolatedNetwork && nic.NetworkID == networkID {
			return nic
		}
	}
	return nil
}

func findNICByNetwork(nics []map[string]interface{}, networkID string) map[string]interface{} {
	for _, nic := range nics {
		if nic["network_id"].(string) == networkID {
			return nic
		}
	}
	return nil
}

func findNICByID(nics []map[string]interface{}, nicID int) map[string]interface{} {
	for _, nic := range nics {
		if nic["id"].(int) == nicID {
			return nic
		}
	}
	return nil
}
