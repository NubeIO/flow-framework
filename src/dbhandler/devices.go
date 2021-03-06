package dbhandler

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

func (h *Handler) GetDevices(args api.Args) ([]*model.Device, error) {
	q, err := getDb().GetDevices(args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetDevice(uuid string, args api.Args) (*model.Device, error) {
	q, err := getDb().GetDevice(uuid, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetDeviceByArgs(args api.Args) (*model.Device, error) {
	return getDb().GetOneDeviceByArgs(args)
}

func (h *Handler) DeviceNameExistsInNetwork(deviceName, networkUUID string) (device *model.Device, existing bool) {
	network, err := getDb().GetNetwork(networkUUID, api.Args{WithDevices: true})
	if err != nil {
		return nil, false
	}
	for _, dev := range network.Devices {
		if dev.Name == deviceName {
			return dev, true
		}
	}
	return nil, false
}

func (h *Handler) CreateDevice(body *model.Device) (*model.Device, error) {
	q, err := getDb().CreateDevice(body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateDevice(uuid string, body *model.Device, fromPlugin bool) (*model.Device, error) {
	q, err := getDb().UpdateDevice(uuid, body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) DeleteDevice(uuid string) (bool, error) {
	_, err := getDb().DeleteDevice(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
