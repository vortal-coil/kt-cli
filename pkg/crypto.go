package pkg

import (
	"errors"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// CryptoInfo is a struct that holds all the necessary information for encryption/decryption
type CryptoInfo struct {
	// EncryptedCryptoKey is the encrypted crypto key received from the server. You should decrypt it with a password
	EncryptedCryptoKey string
	// RawCryptoKey is the decrypted crypto key.
	//You can use it for encryption/decryption.
	//It usually is not provided by the server, you should decrypt EncryptedCryptoKey with a password
	RawCryptoKey string
	// PublicKey is the public key of the user. It is used for encryption and signature verification
	PublicKey string
	// Password is used to decrypt the EncryptedCryptoKey, also it is used as passphrase for the private key
	Password string
}

// IsCryptoReady checks if the CryptoInfo is ready for encryption/decryption.
// It should have the crypto key to be decrypted
func (c *CryptoInfo) IsCryptoReady() bool {
	return c.RawCryptoKey != ""
}

// TryGetReady tries to get the CryptoInfo ready for encryption/decryption.
// It tries to decrypt the key with the password.
func (c *CryptoInfo) TryGetReady(token string, disk string) error {
	if c.IsCryptoReady() {
		return nil
	}

	if c.Password == "" && c.RawCryptoKey == "" {
		// Crypto data is provided, but password and key are empty
		return errors.New("no password or decrypted key provided")
	} else if c.RawCryptoKey == "" {
		// Password is provided, but the key is empty. We need to get and decrypt the key
		crypt, err := GetCryptoInfo(token, disk, c.Password)
		if err != nil {
			return fmt.Errorf("failed to get crypto info: %w", err)
		}

		*c = *crypt
	} else {
		return errors.New("no any data provided")
	}

	return nil
}

// GetCryptoInfo gets the CryptoInfo from the server.
// It decrypts the crypto key if it is encrypted and a password is provided
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

// GetKeyRings gets the public and private key rings from the armored keys
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
