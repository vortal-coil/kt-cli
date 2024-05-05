package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"io"
	"net/http"
)

// DownloadFile downloads a file from the cloud. If the file is encrypted, it will be decrypted using the provided.
// If the file is encrypted and no crypto info provided, it will return an error.
// You need to provide at least your crypto password in CryptoInfo to decrypt the file.
// If no keys are provided, it will try to get the crypto info from the server and decrypt your key with the password.
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

	resp, err := MapToStruct[FileGetByIdResponse](filesList.Result)
	if err != nil {
		return "", 0, err
	}

	if resp.Count == 0 {
		return "", 0, errors.New("file not found or you have not access to it")
	}

	list := resp.List
	if len(list) == 0 {
		return "", 0, errors.New("file not found or you have not access to it")
	}

	fileInfo := list[0]

	name := fileInfo.Name
	encrypted := fileInfo.Encrypted
	mimeType := fileInfo.Mime
	disk := fileInfo.Disk

	// If the file is encrypted and no any crypto info provided, we need to get it
	if encrypted && (cryptoInfo == nil || !cryptoInfo.IsCryptoReady()) {
		if cryptoInfo == nil {
			return "", 0, errors.New("file is encrypted but no crypto info provided")
		}

		if err := cryptoInfo.TryGetReady(token, disk); err != nil {
			return "", 0, fmt.Errorf("failed to decrypt file: %w", err)
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

	downloadResponse, err := MapToStruct[DownloadResponse](downloadRequest.Result)
	if err != nil {
		return "", 0, fmt.Errorf("cannot get download link: %w", err)
	}

	fileUrl := downloadResponse.URL
	if len(fileUrl) == 0 {
		return "", 0, errors.New("file url is empty")
	}

	fileResp, err := http.Get(fileUrl)
	if err != nil {
		return "", 0, err
	}
	defer fileResp.Body.Close()

	if fileResp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("bad response status code: %s", fileResp.Status)
	}

	if encrypted {
		currentLogger("File is encrypted, downloading first")
		// At the moment, we download the file to the buffer and then decrypt it.
		// In the future, we will decrypt the file using the stream
		buf := new(bytes.Buffer)
		numBytes, err = io.Copy(buf, fileResp.Body)

		currentLogger("File downloaded. Decrypting now")
		message := crypto.NewPGPMessage(buf.Bytes())

		_, privateKeyRing, err := GetKeyRings(cryptoInfo.PublicKey, cryptoInfo.RawCryptoKey, []byte(cryptoInfo.Password))
		if err != nil {
			return "", 0, err
		}

		decrypted, err := privateKeyRing.Decrypt(message, nil, 0)
		if err != nil {
			return "", 0, err
		}
		privateKeyRing.ClearPrivateParams()

		currentLogger("File decrypted. Saving now")
		numBytes, err = io.Copy(writer, decrypted.NewReader())
	} else {
		currentLogger("File is not encrypted, downloading as-is")
		numBytes, err = io.Copy(writer, fileResp.Body)
	}

	if err != nil {
		return "", 0, err
	}

	currentLogger("Download is done (%d bytes)", numBytes)
	return name, numBytes, nil
}
