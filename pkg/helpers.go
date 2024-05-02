package pkg

import "fmt"

func GetUserDisk(token string, disk string) (string, *CryptoInfo, error) {
	resp, err := ApiRequest(token, "disks.get", nil)
	if err != nil {
		return "", nil, err
	}

	list := resp.Result["list"].([]interface{})
	if len(list) == 0 {
		return "", nil, fmt.Errorf("users default disk not found")
	}

	var diskInfo map[string]interface{}
	for _, nextDisk := range list {
		nextDiskMap := nextDisk.(map[string]interface{})

		// If the disk is set, we need to find the disk with the desired id
		if disk != "" && nextDiskMap["id"].(string) == disk {
			diskInfo = nextDiskMap
			break
		}

		// If the disk is not set, we need to return the first disk as default
		if disk == "" {
			diskInfo = nextDiskMap
			break
		}
	}

	if diskInfo == nil {
		return "", nil, fmt.Errorf("disk not found")
	}

	id := diskInfo["id"].(string)
	cryptoKey := diskInfo["crypto_key"].(string)
	publicKey := diskInfo["public_key"].(string)

	return id, &CryptoInfo{EncryptedCryptoKey: cryptoKey, PublicKey: publicKey}, nil
}
