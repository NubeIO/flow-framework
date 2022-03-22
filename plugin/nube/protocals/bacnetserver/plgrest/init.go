package plgrest

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// RestClient is used to invoke Form3 Accounts API.
type RestClient struct {
	client      *resty.Client
	ClientToken string
}

// NewNoAuth returns a new instance
func NewNoAuth(address string, port string) *RestClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%s", address, port)
	apiURL := url
	client.SetHostURL(apiURL)
	client.SetError(&Error{})
	client.SetHeader("Content-Type", "application/json")
	return &RestClient{client: client}
}
