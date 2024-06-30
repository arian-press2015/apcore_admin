package httpclient

import (
	"net/http"
)

func NewHTTPClient() *http.Client {
	return &http.Client{}
}
