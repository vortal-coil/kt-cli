package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"io"
	"net/http"
)

func DownloadFile(token string, fileId string, writer io.Writer, cryptoInfo *CryptoInfo) (fileName string, numBytes int64, err error) {
	if fileId == "" {
		return "", 0, errors.New("file id is required")
	}

	filesList, err := ApiRequest(token, "files.getById", map[string]interface{}{"file": fileId})
	if err != nil {
		return "", 0, err
	}
	if filesList.Error.Code != 0 {
		return "", 0, errors.New(filesList.Error.Message)
	}

	rawList, ok := filesList.Result["list"]
	if !ok {
		return "", 0, errors.New("bad file get response")
	}

	list, ok := rawList.([]interface{})
	if !ok {
		return "", 0, errors.New("files list parameter is not a list itself")
	}
	if len(list) == 0 {
		return "", 0, errors.New("file not found or you have not access to it")
	}

	fileInfo := list[0].(map[string]interface{})

	name := fileInfo["name"].(string)
	encrypted := fileInfo["encrypted"].(bool)
	mimeType := fileInfo["mime"].(string)
	disk := fileInfo["disk"].(string)

	// If the file is encrypted and no any crypto info provided, we need to get it
	if encrypted && (cryptoInfo == nil || !cryptoInfo.IsCryptoReady()) {
		if cryptoInfo == nil {
			// No any crypto info provided
			return "", 0, errors.New("file is encrypted and no any crypto cryptoInfo provided")
		} else if cryptoInfo.Password == "" && cryptoInfo.RawCryptoKey == "" {
			// Crypto data is provided, but password and key are empty
			return "", 0, errors.New("file is encrypted and no password or keys provided")
		} else if cryptoInfo.RawCryptoKey == "" {
			// Password is provided, but the key is empty. We need to get and decrypt the key
			cryptoInfo, err = GetCryptoInfo(token, disk, cryptoInfo.Password)
			if err != nil {
				return "", 0, fmt.Errorf("failed to get crypto cryptoInfo: %w", err)
			}
		} else {
			return "", 0, errors.New("file is encrypted and no any data provided")
		}
	}

	currentLogger("Downloading file %s (%s)", name, mimeType)

	downloadRequest, err := ApiRequest(token, "files.download", map[string]interface{}{"file": fileId})
	if err != nil {
		return "", 0, err
	}
	if downloadRequest.Error.Code != 0 {
		return "", 0, errors.New(downloadRequest.Error.Message)
	}

	fileUrl, ok := downloadRequest.Result["url"].(string)
	if !ok {
		return "", 0, errors.New("failed to get file url")
	}

	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("bad response status code: %s", resp.Status)
	}

	if encrypted {
		currentLogger("File is encrypted, downloading first")
		buf := new(bytes.Buffer)
		numBytes, err = io.Copy(buf, resp.Body)

		currentLogger("File downloaded. Decrypting now")
		message := crypto.NewPGPMessage(buf.Bytes())

		publicKeyRing, privateKeyRing, err := GetKeyRings(cryptoInfo.PublicKey, cryptoInfo.RawCryptoKey, []byte(cryptoInfo.Password))
		if err != nil {
			return "", 0, err
		}

		decrypted, err := privateKeyRing.Decrypt(message, publicKeyRing, crypto.GetUnixTime())
		privateKeyRing.ClearPrivateParams()

		currentLogger("File decrypted. Saving now")
		numBytes, err = io.Copy(writer, decrypted.NewReader())
	} else {
		currentLogger("File is not encrypted, downloading as-is")
		numBytes, err = io.Copy(writer, resp.Body)
	}

	if err != nil {
		return "", 0, err
	}

	currentLogger("Download is done (%d bytes)", numBytes)
	return name, numBytes, nil
}
