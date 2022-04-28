package main

import (
	"context"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/nube/rubixio"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/rest"
	log "github.com/sirupsen/logrus"
	"time"
)

const pollName = "polling"

type polling struct {
	enable        bool
	loopDelay     time.Duration
	delayNetworks time.Duration
	delayDevices  time.Duration
	delayPoints   time.Duration
	isRunning     bool
}

var poll poller.Poller

func (inst *Instance) updateInputs(pnt *model.Point, inputs *rubixio.Inputs) {
	if inputs == nil {
		return
	}
	IoNumber := pnt.IoNumber //UI1
	pntType := pnt.IoType    //10k
	found, temp, voltage, current, raw, digital, err := rubixio.GetInputValues(inputs, IoNumber)
	if err != nil {
		return
	}
	var pointValue float64
	if found {
		switch pntType {
		case string(model.IOTypeThermistor10K):
			pointValue = temp
		case string(model.IOTypeVoltageDC):
			pointValue = voltage
		case string(model.IOTypeCurrent):
			pointValue = current
		case string(model.IOTypeDigital):
			pointValue = float64(digital)
		case string(model.IOTypeRAW):
			pointValue = raw
		}
		_, err = inst.pointUpdateValue(pnt.UUID, pointValue)
		if err != nil {
			log.Errorln("rubixio.polling.syncInputs() failed to update point value")
		}
	}

}

func (inst *Instance) syncInputs(dev *model.Device, inputs *rubixio.Inputs) {
	for _, pnt := range dev.Points {
		inst.updateInputs(pnt, inputs)
	}
}

func (inst *Instance) syncOutputs(dev *model.Device) (bulk []rubixio.BulkWrite) {
	uo1 := 0
	uo2 := 0
	uo3 := 0
	uo4 := 0
	uo5 := 0
	uo6 := 0
	do1 := 0
	do2 := 0
	for _, pnt := range dev.Points {
		switch pnt.IoNumber {
		case "UO1":
			uo1 = int(nils.Float64IsNil(pnt.WriteValue))
		case "UO2":
			uo2 = int(nils.Float64IsNil(pnt.WriteValue))
		case "UO3":
			uo3 = int(nils.Float64IsNil(pnt.WriteValue))
		case "UO4":
			uo4 = int(nils.Float64IsNil(pnt.WriteValue))
		case "UO5":
			uo5 = int(nils.Float64IsNil(pnt.WriteValue))
		case "UO6":
			uo2 = int(nils.Float64IsNil(pnt.WriteValue))
		case "DO1":
			do1 = int(nils.Float64IsNil(pnt.WriteValue))
		case "DO2":
			do2 = int(nils.Float64IsNil(pnt.WriteValue))
		}
	}
	//todo add in if point is disabled
	bulk = []rubixio.BulkWrite{
		{
			IONum: "UO1",
			Value: uo1,
		},
		{
			IONum: "UO2",
			Value: uo2,
		},
		{
			IONum: "UO3",
			Value: uo3,
		},
		{
			IONum: "UO4",
			Value: uo4,
		},
		{
			IONum: "UO5",
			Value: uo5,
		},
		{
			IONum: "UO6",
			Value: uo6,
		},
		{
			IONum: "DO1",
			Value: do1,
		},
		{
			IONum: "DO2",
			Value: do2,
		},
	}
	return
}

func (inst *Instance) getInputs() *rubixio.Inputs {

	restService := &rest.Service{}
	restService.Url = inst.config.Ip
	restService.Port = 5001
	restOptions := &rest.Options{}
	restService.Options = restOptions
	restService = rest.New(restService)

	nubeProxy := &rest.NubeProxy{}
	restService.NubeProxy = nubeProxy

	bacnetClient := rubixio.New(restService)
	inputs, res := bacnetClient.GetInputs()
	if res.GetError() != nil {
		log.Errorln("rubixio.polling.getInputs() failed to do rest-api call err:", err)
	}
	return inputs
}

func (inst *Instance) writeOutput(dev *model.Device) {
	restService := &rest.Service{}
	restService.Url = inst.config.Ip
	restService.Port = 5001
	restOptions := &rest.Options{}
	restService.Options = restOptions
	restService = rest.New(restService)
	nubeProxy := &rest.NubeProxy{}
	restService.NubeProxy = nubeProxy
	client := rubixio.New(restService)

	bulk := inst.syncOutputs(dev)
	outs, resp := client.UpdatePointValueBulk(bulk)
	if resp.GetError() != nil || outs == nil {
		log.Errorln("rubixio.polling.writeOutput() failed to do rest-api call err:", err)
		return
	}

}

func (inst *Instance) polling(p polling) error {
	var defaultInterval = time.Duration(inst.config.PollingTimeInMs) * time.Millisecond //default polling is 2.5 sec

	if p.delayNetworks <= 0 {
		p.delayNetworks = defaultInterval
	}
	if p.delayDevices <= 0 {
		p.delayDevices = defaultInterval
	}
	if p.delayPoints <= 0 {
		p.delayPoints = defaultInterval
	}
	if p.enable {
		poll = poller.New()
	}
	var counter float64
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	f := func() (bool, error) {
		nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		if len(nets) == 0 {
			time.Sleep(15000 * time.Millisecond)
			log.Info("rubixio-polling: NO NETWORKS FOUND")
		}
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == inst.pluginUUID {
				log.Infof("rubixio-polling: LOOP COUNT: %v\n", counter)
				counter++
				for _, dev := range net.Devices { //DEVICES
					if err != nil {
						log.Errorf("rubixio-polling: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
					}
					dNet := p.delayNetworks
					time.Sleep(dNet)
					inputs := inst.getInputs()
					inst.syncInputs(dev, inputs)
					inst.writeOutput(dev)
					for _, pnt := range dev.Points { //POINTS
						if nils.BoolIsNil(pnt.IsOutput) {

						} else {
							//inst.syncInputs(pnt, inputs)
						}
					}
				}
			}

		}
		if !p.enable { //TODO to disable of the polling isn't working
			return true, nil
		} else {
			return false, nil
		}
	}
	err := poll.Poll(context.Background(), f)
	if err != nil {
		return err
	}
	return nil
}
