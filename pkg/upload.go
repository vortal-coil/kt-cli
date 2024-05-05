package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"io"
	mime2 "mime"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// UploadFile uploads a file to the cloud.
// If encryption is enabled, it will encrypt the file before uploading.
// The file will be encrypted using the public key provided in the CryptoInfo struct.
// You need public key to encrypt the file.
// If you don't have the public key, you can get it from the server using the GetCryptoInfo function.
func UploadFile(token string, name string, rewriteMime string, disk string, folder string, cryptoInfo *CryptoInfo, reader io.Reader) (fileId string, err error) {
	currentLogger("Uploading file %s", name)

	var cryptoVal string
	var publicRing *crypto.KeyRing
	encrypt := cryptoInfo != nil

	if !encrypt {
		cryptoVal = "0"
		confirm := ScanOrDefault("You are uploading a file without encryption. Continue? (y/n): ", "y")
		if confirm != "y" {
			return "", errors.New("upload aborted by user")
		}
		currentLogger("Uploading without encryption")
	} else {
		cryptoVal = "1"
		currentLogger("Encrypting")

		// If the crypto info is not ready, we need to get it. Nil check is not necessary because we have already checked it
		if !cryptoInfo.IsCryptoReady() {
			if err := cryptoInfo.TryGetReady(token, disk); err != nil {
				return "", fmt.Errorf("failed to encrypt file: %w", err)
			}
		}

		publicRing, _, err = GetKeyRings(cryptoInfo.PublicKey, cryptoInfo.RawCryptoKey, []byte(cryptoInfo.Password))
		if err != nil {
			return "", err
		}
		if !publicRing.CanEncrypt() {
			return "", errors.New("public key cannot encrypt")
		}
	}

	client := &http.Client{}
	body := &bytes.Buffer{}
	writerMultipart := multipart.NewWriter(body)

	_ = writerMultipart.WriteField("token", token)
	_ = writerMultipart.WriteField("disk", strings.TrimSpace(disk))
	_ = writerMultipart.WriteField("folder", strings.TrimSpace(folder))
	_ = writerMultipart.WriteField("crypto", strings.TrimSpace(cryptoVal))
	part, err := writerMultipart.CreateFormFile("file", name)
	if err != nil {
		return "", err
	}

	if encrypt {
		messageMeta := crypto.NewPlainMessageMetadata(true, name, time.Now().Unix())

		plainWriter, err := publicRing.EncryptStreamWithCompression(part, messageMeta, nil)
		if err != nil {
			return "", err
		}

		_, err = io.Copy(plainWriter, reader)
		_ = plainWriter.Close()
	} else {
		_, err = io.Copy(part, reader)
	}

	if err != nil {
		return "", err
	}

	err = writerMultipart.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", uploadUrl, body)
	if err != nil {
		return "", err
	}

	var mime string
	if rewriteMime != "" {
		mime = rewriteMime
	} else {
		mime = mime2.TypeByExtension(name)
		if mime == "" {
			mime = writerMultipart.FormDataContentType()
		}
	}
	req.Header.Set("Content-Type", mime)

	currentLogger("Uploading file to server")
	responseInfo, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer responseInfo.Body.Close()

	rawResponse, err := readerToMap(responseInfo.Body)
	if err != nil {
		return "", err
	}

	response, err := MapToStruct[ApiResponse](rawResponse)
	if err != nil {
		return "", err
	}

	if response.Error.Code != 0 {
		return "", fmt.Errorf("%s: %s (code %d)", responseInfo.Status, response.Error.Message, response.Error.Code)
	} else if responseInfo.StatusCode != http.StatusOK {
		return "", errors.New(responseInfo.Status)
	}

	result, err := MapToStruct[UploadResult](response.Result)
	if err != nil {
		return "", err
	}

	if result.Ok {
		fileId = result.FileID
		if fileId == "" {
			return "", errors.New("response file_id is empty")
		}

		currentLogger("File uploaded successfully. File ID: %s", fileId)
		return fileId, nil
	}

	return "", errors.New("upload failed (unknown reason)")
}
