package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (a *FlowClient) Login(body *model.LoginBody) (*model.Token, error) {
	resp, err := a.client.R().
		SetBody(body).
		SetResult(&model.Token{}).
		Post("/api/users/login")

	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("login: %s", err)
		} else {
			return nil, fmt.Errorf("login: %s", resp)
		}
	}
	if resp.IsError() {
		return nil, fmt.Errorf("login: %s", resp)
	}
	return resp.Result().(*model.Token), nil
}
