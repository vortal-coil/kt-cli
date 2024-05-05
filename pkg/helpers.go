package pkg

import "fmt"

// GetUserDisk returns the user's default disk or the disk with the desired id.
// It also returns the crypto info for the disk
func GetUserDisk(token string, disk string) (*Disk, *CryptoInfo, error) {
	resp, err := ApiRequest(token, "disks.get", nil)
	if err != nil {
		return nil, nil, err
	}

	// At the moment results are not a structure, so we need to cast it to a map
	disks, err := MapToStruct[DisksInfo](resp.Result)
	if err != nil {
		return nil, nil, err
	}
	if len(disks.List) == 0 {
		return nil, nil, fmt.Errorf("users default disk not found")
	}

	var diskInfo *Disk
	for _, nextDisk := range disks.List {
		// If the disk is set, we need to find the disk with the desired id
		if disk != "" && nextDisk.ID == disk {
			diskInfo = nextDisk
			break
		}

		// If the disk is not set, we need to return the first disk as default
		if disk == "" {
			diskInfo = nextDisk
			break
		}
	}

	if diskInfo == nil {
		return nil, nil, fmt.Errorf("disk not found")
	}

	cryptoKey := diskInfo.CryptoKey
	publicKey := diskInfo.PublicKey
	return diskInfo, &CryptoInfo{EncryptedCryptoKey: cryptoKey, PublicKey: publicKey}, nil
}
