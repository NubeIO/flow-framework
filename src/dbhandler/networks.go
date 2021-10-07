package dbhandler

import (
	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetNetworkByPlugin(pluginUUID string, args api.Args) (*model.Network, error) {
	q, err := getDb().GetNetworkByPlugin(pluginUUID, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworksByPlugin(pluginUUID string, args api.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworksByPlugin(pluginUUID, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetNetworksByName(name string, args api.Args) ([]*model.Network, error) {
	q, err := getDb().GetNetworksByName(name, args)
	if err != nil {
		return nil, err
	}
	return q, nil
}
