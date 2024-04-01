package pkg

import (
	"errors"
	"github.com/kt-soft-dev/kt-cli/internal"
)

func CheckToken(token string) (string, error) {
	request, err := ApiRequest(token, "auth.getMe", nil)
	if err != nil {
		return "", err
	}
	if request.Error.Code != 0 {
		internal.Print(request.Error.Message)
		return "", errors.New(request.Error.Message)
	}

	if id, ok := request.Result["id"]; ok {
		return id.(string), nil
	}

	return "", errors.New("failed to get user id from response")
}
