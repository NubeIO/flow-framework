package client

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/rest/v1/rest"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/go-resty/resty/v2"
)

// ClientGetPlugins an object
func (a *FlowClient) ClientGetPlugins() (*ResponsePlugins, error) {
	resp, err := a.client.R().
		SetResult(&ResponsePlugins{}).
		Get("/plugins")
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*ResponsePlugins), nil
}

// ClientGetPlugin an object
func (a *FlowClient) ClientGetPlugin(uuid string) (*ResponseBody, error) {
	resp, err := a.client.R().
		SetResult(&ResponseBody{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/plugins/{uuid}")
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*ResponseBody), nil
}

func failedResponse(err error, resp *resty.Response) error {
	if err != nil {
		return fmt.Errorf("%s failed", err)
	}
	if resp.Error() != nil {
		return getAPIError(resp)
	}
	if rest.StatusCodesAllBad(resp.StatusCode()) {
		return getAPIError(resp)
	}
	return nil

}

// CreateNetworkPlugin an object
func (a *FlowClient) CreateNetworkPlugin(body *model.Network, pluginName string) (*model.Network, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/networks", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Network{}).
		SetBody(body).
		Post(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*model.Network), nil
}

// DeleteNetworkPlugin delete an object
func (a *FlowClient) DeleteNetworkPlugin(body *model.Network, pluginName string) (ok bool, err error) {
	url := fmt.Sprintf("/api/plugins/api/%s/networks", pluginName)
	resp, err := a.client.R().
		SetBody(body).
		Delete(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return false, failResponse
	}
	return true, err
}

// DeleteDevicePlugin delete an object
func (a *FlowClient) DeleteDevicePlugin(body *model.Device, pluginName string) (ok bool, err error) {
	url := fmt.Sprintf("/api/plugins/api/%s/devices", pluginName)
	resp, err := a.client.R().
		SetBody(body).
		Delete(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return false, failResponse
	}
	return true, err
}

// DeletePointPlugin delete an object
func (a *FlowClient) DeletePointPlugin(body *model.Point, pluginName string) (ok bool, err error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points", pluginName)
	resp, err := a.client.R().
		SetBody(body).
		Delete(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return false, failResponse
	}
	return true, err
}

// CreateDevicePlugin an object
func (a *FlowClient) CreateDevicePlugin(body *model.Device, pluginName string) (*model.Device, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/devices", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Device{}).
		SetBody(body).
		Post(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
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
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*model.Point), nil
}

// UpdateNetworkPlugin update an object
func (a *FlowClient) UpdateNetworkPlugin(body *model.Network, pluginName string) (*model.Network, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/networks", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Network{}).
		SetBody(body).
		Patch(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*model.Network), nil
}

// UpdateDevicePlugin update an object
func (a *FlowClient) UpdateDevicePlugin(body *model.Device, pluginName string) (*model.Device, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/devices", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Device{}).
		SetBody(body).
		Patch(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*model.Device), nil
}

// UpdatePointPlugin update an object
func (a *FlowClient) UpdatePointPlugin(body *model.Point, pluginName string) (*model.Point, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		Patch(url)
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*model.Point), nil
}

// WritePointPlugin update an object
func (a *FlowClient) WritePointPlugin(pointUUID string, body *model.Point, pluginName string) (*model.Point, error) {
	url := fmt.Sprintf("/api/plugins/api/%s/points/write/{uuid}", pluginName)
	resp, err := a.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": pointUUID}).
		Patch(url)
	fmt.Println(fmt.Sprintf("WritePointPlugin %+v", resp))
	failResponse := failedResponse(err, resp)
	if failResponse != nil {
		return nil, failResponse
	}
	return resp.Result().(*model.Point), nil
}
