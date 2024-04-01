package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"unicode"
)

// JsonToReader converts a map to a ReadCloser. It makes easier to send json data to a http request
func JsonToReader(data map[string]interface{}) io.ReadCloser {
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

// JsonToString converts a map to a string. It makes easier to print json data. Pretty-prints the json if pretty is true
func JsonToString(data map[string]interface{}, pretty bool) string {
	if pretty {
		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return ""
		}

		return string(jsonData)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

func ParseKeyValues(data string) map[string]interface{} {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	items := strings.FieldsFunc(data, f)
	m := make(map[string]interface{})
	for _, item := range items {
		x := strings.Split(item, "=")
		m[x[0]] = x[1]
	}

	return m
}

// GetActualError returns an error if the response has an error or if there was an error passed as argument
// It is useful to avoid checking if the response is nil and if the response error is nil every time you make a request
func GetActualError(response *ApiResponse, err error) error {
	if err != nil {
		return err
	}
	if response != nil && response.Error.Code != 0 {
		return errors.New(response.Error.Message)
	}

	return nil
}
