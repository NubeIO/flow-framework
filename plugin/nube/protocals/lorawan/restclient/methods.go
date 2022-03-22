package lwrest

import (
	"fmt"
	lwmodel "github.com/NubeIO/flow-framework/plugin/nube/protocals/lorawan/lorawan_model"
)

const limit = "50"
const orgID = "1"

// GetOrganizations get all
func (a *RestClient) GetOrganizations() (*lwmodel.Organizations, error) {
	q := fmt.Sprintf("/api/organizations?limit=%s", limit)
	resp, err := a.client.R().
		SetResult(lwmodel.Organizations{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetOrganizations %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.Organizations), nil
}

// GetGateways get all
func (a *RestClient) GetGateways() (*lwmodel.Gateways, error) {
	q := fmt.Sprintf("/api/gateways?limit=%s&organizationID=%s", limit, orgID)
	resp, err := a.client.R().
		SetResult(lwmodel.Gateways{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetApplications %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.Gateways), nil
}

// GetApplications get all
func (a *RestClient) GetApplications() (*lwmodel.Applications, error) {
	q := fmt.Sprintf("/api/applications?limit=%s&organizationID=%s", limit, orgID)
	resp, err := a.client.R().
		SetResult(lwmodel.Applications{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetApplications %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.Applications), nil
}

// GetDevices get all
func (a *RestClient) GetDevices() (*lwmodel.Devices, error) {
	q := fmt.Sprintf("/api/devices?limit=%s&organizationID=%s", limit, orgID)
	resp, err := a.client.R().
		SetResult(lwmodel.Devices{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetDevices %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.Devices), nil
}

// GetDeviceProfiles get all
func (a *RestClient) GetDeviceProfiles() (*lwmodel.DeviceProfiles, error) {
	q := fmt.Sprintf("/api/device-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := a.client.R().
		SetResult(lwmodel.DeviceProfiles{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetDeviceProfiles %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.DeviceProfiles), nil
}

// GetServiceProfiles get all
func (a *RestClient) GetServiceProfiles() (*lwmodel.ServiceProfiles, error) {
	q := fmt.Sprintf("/api/service-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := a.client.R().
		SetResult(lwmodel.ServiceProfiles{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetServiceProfiles %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.ServiceProfiles), nil
}

// GetGatewayProfiles get all
func (a *RestClient) GetGatewayProfiles() (*lwmodel.GatewayProfiles, error) {
	q := fmt.Sprintf("/api/gateway-profiles?limit=%s&organizationID=%s", limit, orgID)
	resp, err := a.client.R().
		SetResult(lwmodel.GatewayProfiles{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetGatewayProfiles %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.GatewayProfiles), nil
}

// AddDevice add all
func (a *RestClient) AddDevice(body lwmodel.Devices) (*lwmodel.Devices, error) {
	q := fmt.Sprintf("/api/devices")
	resp, err := a.client.R().
		SetResult(lwmodel.Devices{}).
		SetBody(body).
		Post(q)
	if err != nil {
		return nil, fmt.Errorf("AddDevice %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.Devices), nil
}

// GetDevice get an object
func (a *RestClient) GetDevice(devEui string) (*lwmodel.GetDevice, error) {
	q := fmt.Sprintf("/api/devices/%s", devEui)
	resp, err := a.client.R().
		SetResult(lwmodel.GetDevice{}).
		Get(q)
	if err != nil {
		return nil, fmt.Errorf("GetDevice %s failed", err)
	}
	if resp.Error() != nil {
		return nil, getAPIError(resp)
	}
	return resp.Result().(*lwmodel.GetDevice), nil
}

// EditDevice edit object
func (a *RestClient) EditDevice(devEui string, body lwmodel.Device) (bool, error) {
	q := fmt.Sprintf("/api/devices/%s", devEui)
	resp, err := a.client.R().
		SetResult(lwmodel.Device{}).
		SetBody(body).
		Put(q)
	if err != nil {
		return false, fmt.Errorf("EditDevice %s failed", err)
	}
	if resp.Error() != nil {
		return false, getAPIError(resp)
	}
	return true, nil
}

// DeleteDevice delete
func (a *RestClient) DeleteDevice(devEui string) (bool, error) {
	q := fmt.Sprintf("/api/devices/%s", devEui)
	resp, err := a.client.R().
		Delete(q)
	if err != nil {
		return false, fmt.Errorf("DeleteDevice %s failed", err)
	}
	if resp.Error() != nil {
		return false, getAPIError(resp)
	}
	return true, nil
}
