package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kt-soft-dev/kt-cli/pkg"
	"os"
	"strings"
	"unicode"
)

// JsonToString converts a map to a string. It makes it easier to print json data. Pretty-prints the json if pretty is true
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

// ParseKeyValues parses a string with key=value pairs separated by spaces and returns a map with the key and value
// It is used to parse the flags of the command line and other similar cases
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
		if len(x) == 2 {
			m[x[0]] = x[1]
		}
	}

	return m
}

// GetActualError returns an error if the response has an error or if there was an error passed as argument.
// It is useful to avoid checking if the response is nil, and if the response error is nil every time you make a request
func GetActualError(response *pkg.ApiResponse, err error) error {
	if err != nil {
		return err
	}
	if response != nil && response.Error.Code != 0 {
		return errors.New(response.Error.Message)
	}

	return nil
}

// IsStdin checks if the stdin has data
func IsStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		PrintError("os.Stdin.Stat()", err)
		return false
	}

	return fi.Size() > 0
}

// ByteCount converts bytes to human-readable format
func ByteCount(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
