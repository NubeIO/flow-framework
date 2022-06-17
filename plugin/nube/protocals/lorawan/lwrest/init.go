package lwrest

import (
	"fmt"
	"github.com/NubeIO/flow-framework/nresty"
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

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewChirp returns a new instance of NewChirp.
func NewChirp(name string, password string, address string, port string) *RestClient {
	log.Infof("lorawan: try and connect to chirpstack on user:%s password:*** ip:%s, port:%s", name, address, port)
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%s", address, port)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Content-Type", "application/json")
	var t Token
	resp, err := client.R().
		SetResult(&t).
		SetHeader("Content-Type", "application/json").
		SetBody(user{Email: name, Password: password}).
		Post("/api/internal/login")
	if err != nil {
		log.Println("getToken err:", err, resp.Status())
	}
	client.SetHeader("Grpc-Metadata-Authorization", t.JWT)
	return &RestClient{client: client, ClientToken: t.JWT}
}
