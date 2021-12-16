package database

import (
	"errors"
	"fmt"

	"github.com/NubeIO/flow-framework/model"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) WizardNewNetworkDevicePoint(plugin string, network *model.Network, device *model.Device, point *model.Point) (*model.Point, error) {
	if network == nil {
		network = &model.Network{
			TransportType: "ip",
		}
	}
	if device == nil {
		device = &model.Device{}
	}

	p, err := d.GetPluginByPath(plugin)
	if err != nil {
		return nil, errors.New("not valid plugin found")
	}

	network.PluginConfId = p.UUID
	n, err := d.CreateNetwork(network)
	if err != nil {
		return nil, fmt.Errorf("network creation failure: %s", err)
	}
	log.Info("Created a Network")

	device.NetworkUUID = n.UUID
	dev, err := d.CreateDevice(device)
	if err != nil {
		return nil, fmt.Errorf("device creation failure: %s", err)
	}
	log.Info("Created a Device: ", dev)

	if point != nil {
		point.DeviceUUID = dev.UUID
		_, err = d.CreatePoint(point, "")
		if err != nil {
			return nil, fmt.Errorf("consumer point creation failure: %s", err)
		}
		log.Info("Created a Point for Consumer", point)
	}
	return point, nil
}
