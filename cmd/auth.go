package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
)

func loadToken() (string, error) {
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("not logged in, please log in first")
		}
		return "", err
	}
	return string(token), nil
}

func authenticateRequest(req *http.Request) error {
	token, err := loadToken()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}
