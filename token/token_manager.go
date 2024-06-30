package token

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arian-press2015/apcore_admin/config"
)

type TokenManager struct {
	Config *config.Config
}

func NewTokenManager(config *config.Config) *TokenManager {
	return &TokenManager{Config: config}
}

func (tm *TokenManager) SaveToken(token string) error {
	return os.WriteFile(tm.Config.TokenFile, []byte(token), 0600)
}

func (tm *TokenManager) LoadToken() (string, error) {
	token, err := os.ReadFile(tm.Config.TokenFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("not logged in, please log in first")
		}
		return "", err
	}
	return string(token), nil
}

func (tm *TokenManager) AuthenticateRequest(req *http.Request) error {
	token, err := tm.LoadToken()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}
