package pkg

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

type CryptoInfo struct {
	EncryptedCryptoKey string
	RawCryptoKey       string
	PublicKey          string
	Password           string
}

func (c *CryptoInfo) IsCryptoReady() bool {
	return c.RawCryptoKey != ""
}

func GetCryptoInfo(token string, disk string, password string) (*CryptoInfo, error) {
	_, cryptoInfo, err := GetUserDisk(token, disk)
	if err != nil {
		return nil, err
	}

	if password != "" {
		cryptoInfo.Password = password
	}

	// We need to decrypt the crypto key if it is not yet decrypted, and we have a password
	if cryptoInfo.RawCryptoKey == "" && cryptoInfo.EncryptedCryptoKey != "" && password != "" {
		message, err := helper.DecryptMessageWithPassword([]byte(password), cryptoInfo.EncryptedCryptoKey)
		if err != nil {
			return cryptoInfo, err
		}

		cryptoInfo.RawCryptoKey = message
	}

	return cryptoInfo, nil
}

func GetKeyRings(publicKey string, privateKey string, passwd []byte) (public *crypto.KeyRing, private *crypto.KeyRing, err error) {
	privateKeyObj, err := crypto.NewKeyFromArmored(privateKey)
	if err != nil {
		return nil, nil, err
	}

	unlockedKeyObj, err := privateKeyObj.Unlock(passwd)
	if err != nil {
		return nil, nil, err
	}

	private, err = crypto.NewKeyRing(unlockedKeyObj)
	if err != nil {
		return nil, nil, err
	}

	publicKeyObj, err := crypto.NewKeyFromArmored(publicKey)
	if err != nil {
		return nil, nil, err
	}

	public, err = crypto.NewKeyRing(publicKeyObj)
	if err != nil {
		return nil, nil, err
	}

	return public, private, nil
}
