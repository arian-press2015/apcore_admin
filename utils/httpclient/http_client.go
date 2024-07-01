package httpclient

import (
	"fmt"
	"io"
	"net/http"

	"github.com/arian-press2015/apcore_admin/config"
	"github.com/arian-press2015/apcore_admin/token"
)

type HTTPClient struct {
	Client       *http.Client
	TokenManager *token.TokenManager
}

func NewHTTPClient(cfg *config.Config) *HTTPClient {
	return &HTTPClient{
		Client:       &http.Client{},
		TokenManager: token.NewTokenManager(cfg),
	}
}

func (hc *HTTPClient) MakeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	err = hc.TokenManager.AuthenticateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %v", err)
	}

	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error: %s", string(body))
	}

	return resp, nil
}

func (hc *HTTPClient) MakeUnauthenticatedRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error: %s", string(body))
	}

	return resp, nil
}
