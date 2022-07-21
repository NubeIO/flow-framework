package main

import (
	"context"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
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

func delays(networkType string) (deviceDelay, pointDelay time.Duration) {
	deviceDelay = 100 * time.Millisecond
	pointDelay = 100 * time.Millisecond
	if networkType == model.TransType.LoRa {
		deviceDelay = 80 * time.Millisecond
		pointDelay = 6000 * time.Millisecond
	}
	return
}

var poll poller.Poller

func (inst *Instance) polling(p polling) error {
	if p.enable {
		poll = poller.New()
	}
	var counter int
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	log.Infoln("init bacnet server network")
	f := func() (bool, error) {
		nets, _ := inst.db.GetNetworksByPlugin(inst.pluginUUID, arg)
		if len(nets) == 0 {
			time.Sleep(2 * time.Second)
			log.Info("bacnet-server: NO NETWORKS FOUND")
		}
		for _, net := range nets { // NETWORKS
			if !inst.pollingEnabled {
				// break
			}
			if net.UUID != "" && net.PluginConfId == inst.pluginUUID {
				timeStart := time.Now()
				devDelay, pointDelay := delays(net.TransportType)
				counter++
				if boolean.IsFalse(net.Enable) {
					log.Infof("bacnet-server: LOOP NETWORK DISABLED: COUNT %v NAME: %s\n", counter, net.Name)
					continue
				}
				for _, dev := range net.Devices { // DEVICES
					if boolean.IsFalse(net.Enable) {
						log.Infof("bacnet-server-device: DEVICE DISABLED: NAME: %s\n", dev.Name)
						continue
					}
					for _, pnt := range dev.Points { // POINTS
						if boolean.IsFalse(net.Enable) {
							continue
						}
						time.Sleep(devDelay) // DELAY between points
						if pnt.WriteMode == "read_only" || pnt.WriteMode == "" {
							readFloat, err := inst.doReadValue(pnt, net.UUID, dev.UUID)
							if err != nil {
								err = inst.pointUpdateErr(pnt.UUID, err)
								continue
							} else {
								_, err := inst.pointUpdateValue(pnt.UUID, readFloat)
								if err != nil {
									continue
								}
							}
						} else if pnt.WriteMode == "write_only" || pnt.WriteMode == "write_once_then_read" {
							// if poll count = 0 or InSync = false then write
							// if write value == nil then don't write
							var doWrite bool
							rsyncWrite := counter % 10
							if counter <= 1 || boolean.IsFalse(pnt.InSync) || rsyncWrite == 0 {
								doWrite = true
								if rsyncWrite == 0 {
									log.Infoln("bacnet-server-WRITE-SYNC-ON-POLL-COUNT on device:", dev.Name, " point:", pnt.Name)
								} else {
									log.Infoln("bacnet-server-WRITE-SYNC on device:", dev.Name, " point:", pnt.Name, " rsyncWrite:", rsyncWrite)
								}
							}
							if float.IsNil(pnt.WriteValue) {
								doWrite = false
								log.Infoln("bacnet-server-WRITE-SYNC-SKIP as writeValue is nil on device:", dev.Name, " point:", pnt.Name, " rsyncWrite:", rsyncWrite)
							}
							if doWrite {
								err := inst.doWrite(pnt, net.UUID, dev.UUID)
								if err != nil {
									err = inst.pointUpdateErr(pnt.UUID, err)
									continue
								}
								// val := float.NonNil(pnt.WriteValue) //TODO not sure is this should then update the PV of the point
								_, err = inst.pointUpdate(pnt.UUID)
								if err != nil {
									continue
								}
							}
						}
						time.Sleep(pointDelay) // DELAY between points
					}
					timeEnd := time.Now()
					diff := timeEnd.Sub(timeStart)
					out := time.Time{}.Add(diff)
					log.Infof("bacnet-server-poll-loop: NETWORK-NAME:%s POLL-DURATION: %s  POLL-COUNT: %d\n", net.Name, out.Format("15:04:05.000"), counter)
				}
			}
		}
		if !p.enable { // TODO the disable of the polling isn't working
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
