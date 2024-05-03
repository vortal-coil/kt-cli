package pkg

import (
	"bytes"
	"encoding/json"
	"io"
)

// jsonToReader converts a map to a ReadCloser. It makes easier to send json data to a http request
func jsonToReader(data map[string]interface{}) io.ReadCloser {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	reader := bytes.NewReader(jsonData)
	if reader == nil {
		return nil
	}

	stringReadCloser := io.NopCloser(reader)
	if stringReadCloser == nil {
		return nil
	}

	return stringReadCloser
}

func readerToMap(reader io.Reader) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.NewDecoder(reader).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
