package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/auth"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/go-resty/resty/v2"
)

type FlowClient struct {
	client      *resty.Client
	ClientToken string
}

func GetFlowToken(ip string, port int, username string, password string) (*string, error) {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d", ip, port)
	client.SetHostURL(url)
	client.SetError(&Error{})
	cli := &FlowClient{client: client}
	token, err := cli.Login(&model.LoginBody{Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	return &token.AccessToken, nil
}

func NewFlowClientCli(ip *string, port *int, token *string, isMasterSlave *bool, globalUUD string, isFNCreator bool) *FlowClient {
	if utils.IsTrue(isMasterSlave) {
		if isFNCreator {
			return newSlaveToMasterCallSession()
		} else {
			return newMasterToSlaveSession(globalUUD)
		}
	} else {
		return newSessionWithToken(*ip, *port, *token)
	}
}

func newSessionWithToken(ip string, port int, token string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("http://%s:%d/ff", ip, port)
	client.SetHostURL(url)
	client.SetError(&Error{})
	client.SetHeader("Authorization", token)
	return &FlowClient{client: client}
}

func newMasterToSlaveSession(globalUUID string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/slave/%s/ff", "0.0.0.0", conf.Server.RSPort, globalUUID)
	client.SetHostURL(url)
	client.SetError(&Error{})
	client.SetHeader("Authorization", auth.GetRubixServiceInternalToken())
	return &FlowClient{client: client}
}

func newSlaveToMasterCallSession() *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/master/ff", "0.0.0.0", conf.Server.RSPort)
	client.SetHostURL(url)
	client.SetError(&Error{})
	client.SetHeader("Authorization", auth.GetRubixServiceInternalToken())
	return &FlowClient{client: client}
}
