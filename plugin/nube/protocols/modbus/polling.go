package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	pollqueue "github.com/NubeIO/flow-framework/plugin/nube/protocols/modbus/poll-queue"
	"github.com/NubeIO/flow-framework/plugin/nube/protocols/modbus/smod"
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
				log.Infof("modbus: No PollingPoint available in Network %s]n", net.UUID)
				continue
			}
			if pp.FFNetworkUUID != net.UUID {
				log.Info("modbus: PollingPoint FFNetworkUUID does not match the Network UUID\n")
				continue
			}
			//fmt.Println("ModbusPolling() pp")
			//fmt.Printf("%+v\n", pp)

			var devArg api.Args
			dev, err := i.db.GetDevice(pp.FFDeviceUUID, devArg)
			if dev == nil || err != nil {
				log.Errorf("modbus: could not find deviceID: %s\n", pp.FFDeviceUUID)
				continue
			}
			if dev.AddressId <= 0 || dev.AddressId >= 255 {
				log.Errorf("modbus: address is not valid.  modbus addresses must be between 1 and 254\n")
				continue
			}

			pnt, err := i.db.GetPoint(pp.FFPointUUID)
			if pnt == nil || err != nil {
				log.Errorf("modbus: could not find pointID: %s\n", pp.FFPointUUID)
				continue
			}

			log.Infof("MODBUS POLL! : Priority: %d, Network: %s Device: %s Point: %s Device-Add: %d Point-Add: %d Point Type: %s \n", pp.PollPriority, net.UUID, dev.UUID, pnt.UUID, dev.AddressId, pnt.AddressID, pnt.ObjectType)

			//fmt.Println("POLLING COMPLETE CALLBACK")
			pollEndTime := time.Now()
			pollDuration := pollEndTime.Sub(pollStartTime)
			pollTimeSecs := pollDuration.Seconds()
			callback(pp, true, true, pollTimeSecs) //(pm *NetworkPollManager) PollingPointCompleteNotification(pp *PollingPoint, writeSuccess, readSuccess bool)

			/*
				var client Client
				//Setup modbus client with Network and Device details
				if net.TransportType == model.TransType.Serial {
					if net.SerialPort != nil || net.SerialBaudRate != nil || net.SerialDataBits != nil || net.SerialStopBits != nil {
						log.Error("modbus: missing serial connection details\n")
						continue
					}
					client.SerialPort = *net.SerialPort
					client.BaudRate = *net.SerialBaudRate
					client.DataBits = *net.SerialDataBits
					client.StopBits = *net.SerialStopBits
					client.Timeout = time.Duration(*net.SerialTimeout) * time.Second
					err = i.setClient(client, net.UUID, true, true)
					if err != nil {
						log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
						continue
					}
				} else {
					client.Host = dev.CommonIP.Host
					client.Port = utils.PortAsString(dev.CommonIP.Port)
					client.Timeout = time.Duration(*net.SerialTimeout) * time.Second
					err = i.setClient(client, net.UUID, true, false)
					if err != nil {
						log.Errorf("modbus: failed to set client %v %s\n", err, dev.CommonIP.Host)
						continue
					}
				}
				cli := getClient()
				address := dev.AddressId
				err = cli.SetUnitId(uint8(address))
				if err != nil {
					log.Errorf("modbus: failed to vaildate SetUnitId %v %d\n", err, dev.AddressId)
					continue
				}
				var ops Operation
				ops.UnitId = uint8(address)
				pnt, err := i.db.GetPoint(pp.FFPointUUID)
				if pp.FFPointUUID != pnt.UUID {
					log.Errorf("modbus: Polling Point FFPointUUID and FF Point UUID don't match\n")
					continue
				}

				if !isConnected() {
					continue
				}
				a := utils.IntIsNil(pnt.AddressID)
				ops.Addr = uint16(a)
				l := utils.IntIsNil(pnt.AddressLength)
				ops.Length = uint16(l)
				ops.ObjectType = pnt.ObjectType
				ops.Encoding = pnt.ObjectEncoding
				ops.IsHoldingReg = utils.BoolIsNil(pnt.IsOutput) //WHY IS THIS HERE?
				ops.ZeroMode = utils.BoolIsNil(dev.ZeroMode)
				_isWrite := isWrite(ops.ObjectType)
				var _pnt model.Point
				_pnt.UUID = pnt.UUID

				//WRITE OPERATION
				//if _isWrite && (pnt.WritePollRequired {    //WRITE ON FIRST PLUGIN ENABLE
				writeSuccess := false
				readSuccess := false
				if _isWrite && utils.BoolIsNil(pnt.WritePollRequired) {
					//WE GET THE WRITE VALUE FROM THE HIGHEST PRIORITY VALUE.  THE PRESENT VALUE IS ONLY SET BY READ OPERATIONS FOR PROTOCOL POINTS
					if pnt.Priority.GetHighestPriorityValue() != nil {
						ops.WriteValue = utils.Float64IsNil(pnt.Priority.GetHighestPriorityValue())
						log.Infof("modbus: WRITE ObjectType: %s  Addr: %d WriteValue: %v\n", ops.ObjectType, ops.Addr, ops.WriteValue)
						request, err := parseRequest(ops)
						if err != nil {
							log.Errorf("modbus parseRequest (WRITE): failed to read holding/input registers: %v\n", err)
						}
						responseRaw, responseValue, err := networkRequest(cli, request)
						log.Infof("modbus: WRITE POLL RESPONSE: ObjectType: %s  Addr: %d  Value:%d  ARRAY: %v\n", ops.ObjectType, ops.Addr, responseValue, responseRaw)
						if err != nil {
							log.Errorf("modbus networkRequest (WRITE): failed to read holding/input registers: %v\n", err)
						}
						if responseValue == ops.WriteValue {
							_pnt.PresentValue = utils.NewFloat64(responseValue)
							_pnt.InSync = utils.NewTrue()
							writeSuccess = true
							readSuccess = true
							_, err = i.pointUpdate(pnt.UUID, &_pnt)
							cov := utils.Float64IsNil(pnt.COV)
							covEvent, _ := utils.COV(ops.WriteValue, utils.Float64IsNil(pnt.OriginalValue), cov)
							if covEvent {
							}
						}
					} else {
						log.Errorf("modbus: no values in priority array to write\n")
					}
				}
				if utils.BoolIsNil(pnt.ReadPollRequired) && !writeSuccess {
					request, err := parseRequest(ops)
					if err != nil {
						log.Errorf("modbus parseRequest (READ): failed to read holding/input registers: %v\n", err)
					}
					responseRaw, responseValue, err := networkRequest(cli, request)
					log.Infof("modbus: WRITE POLL RESPONSE: ObjectType: %s  Addr: %d  Value:%d  ARRAY: %v\n", ops.ObjectType, ops.Addr, responseValue, responseRaw)
					if err != nil {
						log.Errorf("modbus networkRequest (READ): failed to read holding/input registers: %v\n", err)
					} else {
						readSuccess = true
						_pnt.PresentValue = utils.NewFloat64(responseValue)
						_pnt.InSync = utils.NewTrue()
						_, err = i.pointUpdate(pnt.UUID, &_pnt)
						cov := utils.Float64IsNil(pnt.COV)
						covEvent, _ := utils.COV(ops.WriteValue, utils.Float64IsNil(pnt.OriginalValue), cov)
						if covEvent {
						}
						i.store.Set(pnt.UUID, _pnt, -1) //store point in cache
					}
				}
			*/

			// This callback function triggers the PollManager to evaluate whether the point should be re-added to the PollQueue (Never, Immediately, or after the Poll Rate Delay)
			//writeSuccess, readSuccess := true, true
			//callback(pp, writeSuccess, readSuccess)
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
								_, responseValue, err := networkRequest(mbClient, pnt, true)
								if err != nil {
									_, err = i.pointUpdateErr(pnt.UUID, err)
									continue
								}
								_, err = i.pointUpdate(pnt.UUID, responseValue)
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
								_, err = i.pointUpdate(pnt.UUID, responseValue)
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
