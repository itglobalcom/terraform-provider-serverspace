package serverspace

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itglobalcom/terraform-provider-serverspace/serverspace/ssclient"
)

func updateNICS(d *schema.ResourceData, client *ssclient.SSClient, serverID string) error {
	// Preparing network data
	oldNICSValueIfaces, newNICSValueIfaces := d.GetChange("nic")

	oldNICS := convertNICSSetToMap(oldNICSValueIfaces.(*schema.Set))
	newNICS := convertNICSSetToMap(newNICSValueIfaces.(*schema.Set))

	existNICSigns := make([]string, 0)
	for k := range oldNICS {
		existNICSigns = append(existNICSigns, k)

		// exclude non-changed networks
		for _, sign := range existNICSigns {
			if _, ok := newNICS[sign]; ok {
				delete(oldNICS, sign)
				delete(newNICS, sign)
			}
		}
	}

	publicOldNICS := make([]map[string]interface{}, 0)
	privateOldNICS := make([]map[string]interface{}, 0)

	for _, nic := range oldNICS {
		nicNetworkType := ssclient.NetworkType(nic["network_type"].(string))

		if nicNetworkType == ssclient.PublicSharedNetwork {
			publicOldNICS = append(publicOldNICS, nic)
		} else {
			privateOldNICS = append(privateOldNICS, nic)
		}
	}

	publicNewNICS := make([]map[string]interface{}, 0)
	privateNewNICS := make([]map[string]interface{}, 0)

	for _, nic := range newNICS {
		nicNetworkType := ssclient.NetworkType(nic["network_type"].(string))

		if nicNetworkType == ssclient.PublicSharedNetwork {
			publicNewNICS = append(publicNewNICS, nic)
		} else {
			privateNewNICS = append(privateNewNICS, nic)
		}
	}

	if err := updatePublicNICS(client, serverID, publicOldNICS, publicNewNICS); err != nil {
		return err
	}

	return updatePrivateNICS(client, serverID, privateOldNICS, privateNewNICS)
}

func updatePublicNICS(
	client *ssclient.SSClient,
	serverID string,
	oldNICS, newNICS []map[string]interface{},
) error {
	// update exist nics
	processedOldNICCount := 0

	for _, oldNIC := range oldNICS {

		if len(newNICS) == 0 {
			break
		}

		newNIC := newNICS[0]
		newNICS = newNICS[1:]

		nicID := oldNIC["id"].(int)
		newBandwidth := newNIC["bandwidth"].(int)
		if _, err := client.UpdatePublicNICAndWait(serverID, nicID, newBandwidth); err != nil {
			return err
		}
		processedOldNICCount++
	}
	oldNICS = oldNICS[processedOldNICCount:] // remove already updated nics

	for _, oldNIC := range oldNICS {
		if err := client.DeleteNIC(serverID, oldNIC["id"].(int)); err != nil {
			return err
		}
	}

	for _, newNIC := range newNICS {
		newNICBandwidth := newNIC["bandwidth"].(int)
		if _, err := client.CreateNICAndWait(serverID, "", newNICBandwidth); err != nil {
			return err
		}
	}

	return nil
}

func updatePrivateNICS(
	client *ssclient.SSClient,
	serverID string,
	oldNICS, newNICS []map[string]interface{},
) error {
	// update exist nics
	processedOldNICCount := 0

	for _, oldNIC := range oldNICS {

		if len(newNICS) == 0 {
			break
		}

		newNIC := newNICS[0]
		newNICS = newNICS[1:]

		nicID := oldNIC["id"].(int)
		newNetwork := newNIC["network"].(string)
		if err := client.DeleteNIC(serverID, nicID); err != nil {
			return err
		}

		if _, err := client.CreateNICAndWait(serverID, newNetwork, 0); err != nil {
			return err
		}

		processedOldNICCount++
	}
	oldNICS = oldNICS[processedOldNICCount:] // remove already updated nics

	for _, oldNIC := range oldNICS {
		if err := client.DeleteNIC(serverID, oldNIC["id"].(int)); err != nil {
			return err
		}
	}

	for _, newNIC := range newNICS {
		newNetwork := newNIC["network"].(string)
		if _, err := client.CreateNICAndWait(serverID, newNetwork, 0); err != nil {
			return err
		}
	}

	return nil
}

func convertNICSSetToMap(nics *schema.Set) map[string]map[string]interface{} {
	mappedNICS := convertSetToMap(nics)
	signedNICS := make(map[string]map[string]interface{})

	for _, nic := range mappedNICS {
		signedNICS[getNICSignature(nic)] = nic
	}

	return signedNICS
}

func convertSetToMap(set *schema.Set) []map[string]interface{} {
	mapped := make([]map[string]interface{}, set.Len())

	for i, nic := range set.List() {
		mapped[i] = nic.(map[string]interface{})
	}

	return mapped
}

func getNICSignature(nic map[string]interface{}) string {
	netType := nic["network_type"].(string)
	network := nic["network"].(string)
	bandwidth := nic["bandwidth"].(int)
	return fmt.Sprintf("%s-%s-%d", netType, network, bandwidth)
}
