package client

import (
	"fmt"
	"github.com/NubeIO/flow-framework/model"
)

// ClientGetPlugins an object
func (a *FlowClient) ClientGetPlugins() (*ResponsePlugins, error) {
	resp, err := a.client.R().
		SetResult(&ResponsePlugins{}).
		Get("/plugins")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponsePlugins), nil
}

// ClientGetPlugin an object
func (a *FlowClient) ClientGetPlugin(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/plugins/{uuid}")
	if err != nil {
		return nil, fmt.Errorf("fetch name for name %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*ResponseBody), nil
}

// CreateNetworkPlugin an object
func (a *FlowClient) CreateNetworkPlugin(body *model.Network, pluginName string) (*model.Network, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/networks", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Network{}).
		SetBody(body).
		Post(url)
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("add-network: %s", err)
		} else {
			return nil, fmt.Errorf("add-network: %s", resp)
		}
	}
	return resp.Result().(*model.Network), nil
}

// CreateDevicePlugin an object
func (a *FlowClient) CreateDevicePlugin(body *model.Device, pluginName string) (*model.Device, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/devices", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Device{}).
		SetBody(body).
		Post(url)
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("add-device: %s", err)
		} else {
			return nil, fmt.Errorf("add-device: %s", resp)
		}
	}
	return resp.Result().(*model.Device), nil
}

// CreatePointPlugin an object
func (a *FlowClient) CreatePointPlugin(body *model.Point, pluginName string) (*model.Point, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		Post(url)
	if err != nil {
		if resp == nil || resp.String() == "" {
			return nil, fmt.Errorf("add-point: %s", err)
		} else {
			return nil, fmt.Errorf("add-point: %s", resp)
		}
	}
	return resp.Result().(*model.Point), nil
}
