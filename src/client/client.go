package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/auth"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/nresty"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/go-resty/resty/v2"
)

type FlowClient struct {
	client      *resty.Client
	ClientToken string
}

func GetFlowToken(ip string, port int, username string, password string) (*string, error) {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("%s://%s:%d", getSchema(port), ip, port)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	cli := &FlowClient{client: client}
	token, err := cli.Login(&model.LoginBody{Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	return &token.AccessToken, nil
}

func NewLocalClient() (cli *FlowClient) {
	client := resty.New()
	client.SetDebug(false)
	port := config.Get().Server.Port
	url := fmt.Sprintf("%s://%s:%d", getSchema(port), "0.0.0.0", port)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	cli = &FlowClient{client: client}
	return
}

func NewFlowClientCliFromFN(fn *model.FlowNetwork) *FlowClient {
	if boolean.IsTrue(fn.IsMasterSlave) {
		return newSlaveToMasterCallSession()
	} else {
		if boolean.IsTrue(fn.IsRemote) {
			return newSessionWithToken(*fn.FlowIP, *fn.FlowPort, *fn.FlowToken)
		} else {
			return NewLocalClient()
		}
	}
}

func NewFlowClientCliFromFNC(fnc *model.FlowNetworkClone) *FlowClient {
	if boolean.IsTrue(fnc.IsMasterSlave) {
		return newMasterToSlaveSession(fnc.GlobalUUID)
	} else {
		if boolean.IsTrue(fnc.IsRemote) {
			return newSessionWithToken(*fnc.FlowIP, *fnc.FlowPort, *fnc.FlowToken)
		} else {
			return NewLocalClient()
		}
	}
}

func newSessionWithToken(ip string, port int, token string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	url := fmt.Sprintf("%s://%s:%d/ff", getSchema(port), ip, port)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", token)
	return &FlowClient{client: client}
}

func newMasterToSlaveSession(globalUUID string) *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/slave/%s/ff", "0.0.0.0", conf.Server.RSPort, globalUUID)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", auth.GetRubixServiceInternalToken())
	return &FlowClient{client: client}
}

func newSlaveToMasterCallSession() *FlowClient {
	client := resty.New()
	client.SetDebug(false)
	conf := config.Get()
	url := fmt.Sprintf("http://%s:%d/master/ff", "0.0.0.0", conf.Server.RSPort)
	client.SetBaseURL(url)
	client.SetError(&nresty.Error{})
	client.SetHeader("Authorization", auth.GetRubixServiceInternalToken())
	return &FlowClient{client: client}
}

func getSchema(port int) string {
	if port == 443 {
		return "https"
	}
	return "http"
}
