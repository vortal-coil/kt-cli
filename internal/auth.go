package internal

import "github.com/kt-soft-dev/kt-cli/pkg"

func CheckTokenAndAssign(token string, config *Config) error {
	id, err := CheckToken(token)
	if err != nil {
		PrintError(err.Error())
		return nil
	}

	Print("Logged in as user id %s", id)
	config.Token = token
	config.UserID = id
	return nil
}

func CheckToken(token string) (string, error) {
	id, err := pkg.CheckToken(token)
	if err != nil {
		PrintError("Failed to check token")
		return "", err
	}

	return id, nil
}
