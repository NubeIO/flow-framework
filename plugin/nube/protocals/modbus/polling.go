package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/src/poller"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
	log "github.com/sirupsen/logrus"
)

const defaultInterval = 1000 * time.Millisecond

func checkDevValid(uuid string) (bool, error) {
	if uuid == "" {
		log.Errorf("modbus: device id is null \n")
		return false, errors.New("modbus: failed to set client")
	}
	return true, nil
}

func valueRaw(responseRaw interface{}) []byte {
	j, err := json.Marshal(responseRaw)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	return j
}

//TODO: currently Polling loops through each network, grabs one point, and polls it.  Could be improved by having a seperate client/go routine for each of the networks.
func (i *Instance) ModbusPolling() error {
	var poll poller.Poller
	var counter int
	f := func() (bool, error) {
		var netArg api.Args
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, netArg)
		if err != nil {
			return false, err
		}
		if len(nets) == 0 {
			time.Sleep(15000 * time.Millisecond) //WHAT DOES THIS LINE DO?
			log.Info("modbus: NO MODBUS NETWORKS FOUND\n")
		}
		for _, net := range nets { //LOOP THROUGH AND POLL NEXT POINTS IN EACH NETWORK QUEUE
			if net.UUID != "" && net.PluginConfId == i.pluginUUID && net.PollManager != nil {
				log.Infof("modbus: LOOP COUNT: %v\n", counter)
				counter++

				pp, callback := net.PollManager.GetNextPollingPoint() //TODO: once polling completes, callback should be called
				if pp == nil {
					log.Infof("modbus: No PollingPoint available in Network %s]n", net.UUID)
					continue
				}
				if pp.FFNetworkUUID != net.UUID {
					log.Info("modbus: PollingPoint FFNetworkUUID does not match the Network UUID\n")
					continue
				}
				validDev, err := checkDevValid(pp.FFDeviceUUID)
				if !validDev || err != nil {
					continue
				}
				var devArg api.Args
				dev, err := i.db.GetDevice(pp.FFDeviceUUID, devArg)
				if dev.AddressId <= 0 || dev.AddressId >= 255 {
					log.Errorf("modbus: address is not valid.  modbus addresses must be between 1 and 254\n")
					continue
				}

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
				pnt, err := i.db.GetPoint(pp.FFDeviceUUID)
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
				// This callback function triggers the PollManager to evaluate whether the point should be re-added to the PollQueue (Never, Immediately, or after the Poll Rate Delay)
				callback(pp, writeSuccess, readSuccess)
			}
		}
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	i.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}
