package lwrest

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type RestClient struct {
	client      *resty.Client
	ClientToken string
}

type Token struct {
	JWT string `json:"jwt"`
}

// NewChirp returns a new instance of NewChirp.
func NewChirp(name string, password string, address string, port string) *RestClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%s", address, port)
	apiURL := url
	client.SetHostURL(apiURL)
	client.SetError(&Error{})
	client.SetHeader("Content-Type", "application/json")
	var t Token
	getToken, err := client.R().
		SetResult(&t).
		SetHeader("Content-Type", "application/json").
		SetBody(`{"email":"admin", "password":"admin"}`).
		Post("/api/internal/login")
	if err != nil {
		log.Println("getToken err:", err, getToken.Status())
	}
	client.SetHeader("Grpc-Metadata-Authorization", t.JWT)
	return &RestClient{client: client, ClientToken: t.JWT}
}
