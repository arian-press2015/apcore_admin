package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/arian-press2015/apcore_admin/config"
)

type HTTPParser struct {
	Client *HTTPClient
}

func NewHTTPParser(cfg *config.Config) *HTTPParser {
	return &HTTPParser{
		Client: NewHTTPClient(cfg),
	}
}

func (p *HTTPParser) ParseRequest(method string, url string, body interface{}, response interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := p.Client.MakeRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	return nil
}


func (p *HTTPParser) ParseUnauthenticatedRequest(method string, url string, body interface{}, response interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := p.Client.MakeUnauthenticatedRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	return nil
}
