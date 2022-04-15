package main

import (
	"context"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

type polling struct {
	enable        bool
	loopDelay     time.Duration
	delayNetworks time.Duration
	delayDevices  time.Duration
	delayPoints   time.Duration
	isRunning     bool
}

type devCheck struct {
	devUUID string
	client  Client
}

func delays(networkType string) (deviceDelay, pointDelay time.Duration) {
	deviceDelay = 250 * time.Millisecond
	pointDelay = 500 * time.Millisecond
	if networkType == model.TransType.LoRa {
		deviceDelay = 80 * time.Millisecond
		pointDelay = 6000 * time.Millisecond
	}
	return
}

var poll poller.Poller

func (inst *Instance) PollingTCP(p polling) error {
	if p.enable {
		poll = poller.New()
	}
	var counter int
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	f := func() (bool, error) {
		nets, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, arg)
		if len(nets) == 0 {
			time.Sleep(2 * time.Second)
			log.Info("modbus: NO MODBUS NETWORKS FOUND")
		}

		for _, net := range nets { //NETWORKS
			if !inst.pollingEnabled {
				break
			}
			if net.UUID != "" && net.PluginConfId == inst.pluginUUID {
				timeStart := time.Now()
				counter++
				log.Infof("modbus-poll: POLL START: NAME: %s\n", net.Name)
				if !utils.BoolIsNil(net.Enable) {
					log.Infof("modbus: LOOP NETWORK DISABLED: COUNT %v NAME: %s\n", counter, net.Name)
					continue
				}
				for _, dev := range net.Devices { //DEVICES
					if !utils.BoolIsNil(dev.Enable) {
						log.Infof("modbus-device: DEVICE DISABLED: NAME: %s\n", dev.Name)
						continue
					}
					for _, pnt := range dev.Points { //POINTS
						time.Sleep(10 * time.Second)
						if !utils.BoolIsNil(pnt.Enable) {
							log.Infof("modbus-point: POINT DISABLED: NAME: %s\n", pnt.Name)
							continue
						}
						write := isWrite(pnt.ObjectType)
						if write { //IS WRITE
							responseValue := utils.Float64IsNil(pnt.WriteValueOriginal) //feedback in the WriteValue
							_, err = inst.pointUpdate(pnt, responseValue)
						} else { //READ
							responseValue := float64(counter)
							_, err = inst.pointUpdate(pnt, responseValue)
						}
					}
					timeEnd := time.Now()
					diff := timeEnd.Sub(timeStart)
					out := time.Time{}.Add(diff)
					log.Infof("modbus-poll-loop: NETWORK-NAME:%s POLL-DURATION: %s  POLL-COUNT: %d\n", net.Name, out.Format("15:04:05.000"), counter)
				}
			}
		}
		if !p.enable { //TODO the disable of the polling isn't working
			return true, nil
		} else {
			return false, nil
		}
	}
	err := poll.Poll(context.Background(), f)
	if err != nil {
		return nil
	}
	return nil
}
