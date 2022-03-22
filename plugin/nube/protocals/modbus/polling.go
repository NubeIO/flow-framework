package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	pollqueue "github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/poll-queue"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/modbus/smod"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uurl"
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
	deviceDelay = 80 * time.Millisecond
	pointDelay = 80 * time.Millisecond
	if networkType == model.TransType.LoRa {
		deviceDelay = 80 * time.Millisecond
		pointDelay = 6000 * time.Millisecond
	}
	return
}

func (i *Instance) getNetworkPollManagerByUUID(netUUID string) (*pollqueue.NetworkPollManager, error) {
	for _, netPollMan := range i.NetworkPollManagers {
		if netPollMan.FFNetworkUUID == netUUID {
			return netPollMan, nil
		}
	}
	return nil, errors.New("modbus getNetworkPollManagerByUUID(): couldn't find NetworkPollManager")
}

var poll poller.Poller

//TODO: currently Polling loops through each network, grabs one point, and polls it.  Could be improved by having a seperate client/go routine for each of the networks.
func (i *Instance) ModbusPolling() error {
	poll = poller.New()
	var counter = 0
	f := func() (bool, error) {
		counter++
		fmt.Println("\n \n")
		log.Infof("modbus: LOOP COUNT: %v\n", counter)
		var netArg api.Args
		/*
			nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, netArg)
			if err != nil {
				return false, err
			}
		*/

		if len(i.NetworkPollManagers) == 0 {
			//time.Sleep(15000 * time.Millisecond) //WHAT DOES THIS LINE DO?
			log.Info("modbus: NO MODBUS NETWORKS FOUND\n")
		}
		//fmt.Println("i.NetworkPollManagers")
		//fmt.Printf("%+v\n", i.NetworkPollManagers)
		for _, netPollMan := range i.NetworkPollManagers { //LOOP THROUGH AND POLL NEXT POINTS IN EACH NETWORK QUEUE
			log.Infof("ModbusPolling: netPollMan %s", netPollMan.FFNetworkUUID)
			pollStartTime := time.Now()
			//Check that network exists
			//fmt.Println("netPollMan")
			//fmt.Printf("%+v\n", netPollMan)
			net, err := i.db.GetNetwork(netPollMan.FFNetworkUUID, netArg)
			//fmt.Println("net")
			//fmt.Printf("%+v\n", net)
			//fmt.Println("err")
			//fmt.Printf("%+v\n", err)
			if err != nil || net == nil || net.PluginConfId != i.pluginUUID {
				log.Error("modbus: MODBUS NETWORK NOT FOUND\n")
				continue
			}
			//log.Infof("modbus-poll: POLL START: NAME: %s\n", net.Name)

			if !utils.BoolIsNil(net.Enable) {
				log.Infof("modbus: NETWORK DISABLED: COUNT %v NAME: %s\n", counter, net.Name)
				continue
			}
			netPollMan.PrintPollQueuePointUUIDs()
			//fmt.Println("ModbusPolling() current QueueUnloader")
			//fmt.Printf("%+v\n", netPollMan.PluginQueueUnloader.NextPollPoint)
			pp, callback := netPollMan.GetNextPollingPoint() //callback function is called once polling is completed.
			//pp, _ := netPollMan.GetNextPollingPoint() //TODO: once polling completes, callback should be called
			if pp == nil {
				log.Infof("modbus: No PollingPoint available in Network %s", net.UUID)
				continue
			}
			if pp.FFNetworkUUID != net.UUID {
				log.Info("modbus: PollingPoint FFNetworkUUID does not match the Network UUID\n")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			//fmt.Println("ModbusPolling() pp")
			//fmt.Printf("%+v\n", pp)

			var devArg api.Args
			dev, err := i.db.GetDevice(pp.FFDeviceUUID, devArg)
			if dev == nil || err != nil {
				log.Errorf("modbus: could not find deviceID: %s\n", pp.FFDeviceUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			if !utils.BoolIsNil(dev.Enable) {
				log.Errorf("modbus: device is disabled.\n")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			if dev.AddressId <= 0 || dev.AddressId >= 255 {
				log.Errorf("modbus: address is not valid.  modbus addresses must be between 1 and 254\n")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			pnt, err := i.db.GetPoint(pp.FFPointUUID, api.Args{WithPriority: true})
			if pnt == nil || err != nil {
				log.Errorf("modbus: could not find pointID: %s\n", pp.FFPointUUID)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			fmt.Println("ModbusPolling: point")
			fmt.Printf("%+v\n", pnt)

			if pnt.Priority == nil {
				fmt.Println("ModbusPolling: HAD TO ADD PRIORITY ARRAY")
				pnt.Priority = &model.Priority{}
			} else {
				fmt.Println("ModbusPolling: point PRIORITY")
				fmt.Printf("%+v\n", pnt.Priority)
			}

			if !utils.BoolIsNil(pnt.Enable) {
				log.Errorf("modbus: point is disabled.\n")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			log.Infof("MODBUS POLL! : Priority: %s, Network: %s Device: %s Point: %s Device-Add: %d Point-Add: %d Point Type: %s, WriteRequired: %t, ReadRequired: %t \n", pp.PollPriority, net.UUID, dev.UUID, pnt.UUID, dev.AddressId, *pnt.AddressID, pnt.ObjectType, utils.BoolIsNil(pnt.WritePollRequired), utils.BoolIsNil(pnt.ReadPollRequired))

			if !utils.BoolIsNil(pnt.WritePollRequired) && !utils.BoolIsNil(pnt.ReadPollRequired) {
				fmt.Println("polling not required on this point")
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			SetPriorityArrayModeBasedOnWriteMode(pnt) //ensures the point PointPriorityArrayMode is set correctly

			// SETUP MODBUS CLIENT CONNECTION
			var mbClient smod.ModbusClient
			//var dCheck devCheck
			//dCheck.devUUID = dev.UUID
			mbClient, err = i.setClient(net, dev, true)
			if err != nil {
				log.Errorf("modbus: failed to set client error: %v network name:%s\n", err, net.Name)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}
			if net.TransportType == model.TransType.Serial || net.TransportType == model.TransType.LoRa {
				if dev.AddressId >= 1 {
					mbClient.RTUClientHandler.SlaveID = byte(dev.AddressId)
				}
			} else if dev.TransportType == model.TransType.IP {
				url, err := uurl.JoinIpPort(dev.Host, dev.Port)
				if err != nil {
					log.Errorf("modbus: failed to validate device IP %s\n", url)
					netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
					continue
				}
				mbClient.TCPClientHandler.Address = url
				mbClient.TCPClientHandler.SlaveID = byte(dev.AddressId)
			} else {
				log.Errorf("modbus: failed to validate device and network %v %s\n", err, dev.Name)
				netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
				continue
			}

			var responseValue float64
			var response interface{}
			var writeValuePointer *float64
			writeSuccess := false
			if utils.BoolIsNil(pnt.WritePollRequired) { //DO WRITE IF REQUIRED
				log.Info("modbus write point:")
				fmt.Printf("%+v\n", pnt)
				pnt.PrintPointValues()
				writeValuePointer = pnt.Priority.GetHighestPriorityValue()
				if writeValuePointer != nil {
					response, responseValue, err = networkWrite(mbClient, pnt)
					if err != nil {
						_, err = i.pointUpdateErr(pnt.UUID, err)
						netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
						continue
					}
					writeSuccess = true
					log.Infof("modbus-write response: responseValue %f, point UUID: %s, response: %+v \n", responseValue, pnt.UUID, response)
				} else {
					writeSuccess = true //successful because there is no value to write.  Otherwise the point will short cycle.
					log.Info("modbus write point error: no value in priority array to write")
				}
			}

			readSuccess := false
			if utils.BoolIsNil(pnt.ReadPollRequired) { //DO READ IF REQUIRED
				response, responseValue, err = networkRead(mbClient, pnt)
				if err != nil {
					_, err = i.pointUpdateErr(pnt.UUID, err)
					netPollMan.PollingFinished(pp, pollStartTime, false, false, callback)
					continue
				}
				//check cov
				isChange := !utils.CompareFloatPtr(pnt.PresentValue, &responseValue)
				if isChange {
					if err != nil {
						netPollMan.PollingFinished(pp, pollStartTime, writeSuccess, readSuccess, callback)
						continue
					}
				}
				readSuccess = true
				log.Infof("modbus-read response: responseValue %f, point UUID: %s, response: %+v \n", responseValue, pnt.UUID, response)
			}

			//update point in DB if required
			//For write_once and write_always type, write value should become present value
			writeValueToPresentVal := (pnt.WriteMode == poller.WriteOnce || pnt.WriteMode == poller.WriteAlways) && writeSuccess && writeValuePointer != nil

			fmt.Println("ModbusPolling: writeValueToPresentVal ", writeValueToPresentVal)

			if readSuccess || writeValueToPresentVal {
				if writeValueToPresentVal {
					responseValue = *writeValuePointer
					fmt.Println("ModbusPolling: writeOnceWriteValueToPresentVal responseValue: ", responseValue)
					readSuccess = true
				}
				_, err = i.pointUpdate(pnt, responseValue, writeSuccess, readSuccess)
			}

			/*
				//JUST FOR TESTING
				pnt, err = i.db.GetPoint(pp.FFPointUUID)
				if pnt == nil || err != nil {
					log.Errorf("modbus: AFTER... could not find pointID : %s\n", pp.FFPointUUID)
				}
			*/

			// This callback function triggers the PollManager to evaluate whether the point should be re-added to the PollQueue (Never, Immediately, or after the Poll Rate Delay)
			netPollMan.PollingFinished(pp, pollStartTime, writeSuccess, readSuccess, callback)

		}
		time.Sleep(2 * time.Second)
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	i.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}

func (i *Instance) PollingTCP(p polling) error {
	if p.enable {
		poll = poller.New()
	}
	var counter int
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	f := func() (bool, error) {
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if len(nets) == 0 {
			time.Sleep(2 * time.Second)
			log.Info("modbus: NO MODBUS NETWORKS FOUND")
		}

		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				timeStart := time.Now()
				deviceDelay, pointDelay := delays(net.TransportType)
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
					var mbClient smod.ModbusClient
					var dCheck devCheck
					dCheck.devUUID = dev.UUID
					mbClient, err = i.setClient(net, dev, true)
					if err != nil {
						log.Errorf("modbus: failed to set client error: %v network name:%s\n", err, net.Name)
						continue
					}
					if net.TransportType == model.TransType.Serial || net.TransportType == model.TransType.LoRa {
						if dev.AddressId >= 1 {
							mbClient.RTUClientHandler.SlaveID = byte(dev.AddressId)
						}
					} else if dev.TransportType == model.TransType.IP {
						url, err := uurl.JoinIpPort(dev.Host, dev.Port)
						if err != nil {
							log.Errorf("modbus: failed to validate device IP %s\n", url)
							continue
						}
						mbClient.TCPClientHandler.Address = url
						mbClient.TCPClientHandler.SlaveID = byte(dev.AddressId)
					} else {
						log.Errorf("modbus: failed to validate device and network %v %s\n", err, dev.Name)
						continue
					}
					time.Sleep(deviceDelay)          //DELAY between devices
					for _, pnt := range dev.Points { //POINTS
						if !utils.BoolIsNil(pnt.Enable) {
							log.Infof("modbus-point: POINT DISABLED: NAME: %s\n", pnt.Name)
							continue
						}
						write := isWrite(pnt.ObjectType)
						skipDelay := false
						if write { //IS WRITE
							//get existing
							if !utils.BoolIsNil(pnt.InSync) {
								//_, responseValue, err := networkRequest(mbClient, pnt, true)
								if err != nil {
									_, err = i.pointUpdateErr(pnt.UUID, err)
									continue
								}
								//_, err = i.pointUpdate(pnt.UUID, pnt.PointPriorityArrayMode, responseValue)
							} else {
								skipDelay = true
							}
						} else { //READ
							_, responseValue, err := networkRequest(mbClient, pnt, false)
							if err != nil {
								_, err = i.pointUpdateErr(pnt.UUID, err)
								continue
							}
							//simple cov
							isChange := !utils.CompareFloatPtr(pnt.PresentValue, &responseValue)
							if isChange {
								//_, err = i.pointUpdate(pnt.UUID, pnt.PointPriorityArrayMode, responseValue)
								if err != nil {
									continue
								}
							}
						}
						if !skipDelay {
							time.Sleep(pointDelay) //DELAY between points
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
