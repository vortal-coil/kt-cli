package internal

import "github.com/kt-soft-dev/kt-cli/pkg"

// CheckTokenAndAssign checks if the token is valid and assigns it to the config.
// It's a wrapper around CheckToken
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

// CheckToken checks if the token is valid and returns the user id
// It's a wrapper around GetUserID
func CheckToken(token string) (string, error) {
	id, err := pkg.GetUserID(token)
	if err != nil {
		PrintError("Failed to check token")
		return "", err
	}

	return id, nil
}
