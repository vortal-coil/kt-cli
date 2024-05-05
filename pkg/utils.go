package pkg

import (
	"bytes"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
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

// readerToMap converts a reader to a map. It provides an easier way to read json data from http response
func readerToMap(reader io.Reader) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.NewDecoder(reader).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// MapToStruct converts a map to a struct. It is useful when we need to convert a json response to a struct
func MapToStruct[Object any](m map[string]interface{}) (*Object, error) {
	var result Object
	err := mapstructure.Decode(m, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
