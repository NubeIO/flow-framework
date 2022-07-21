package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/edgerest"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/nmath"
	"github.com/NubeIO/flow-framework/utils/structs"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
	"time"
)

const defaultInterval = 10000 * time.Millisecond // default polling is 2.5 sec
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
var getUI *edgerest.UI
var getDI *edgerest.DI

func (inst *Instance) Edge28Polling() error {
	poll = poller.New()
	var counter = 0

	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true

	f := func() (bool, error) {
		counter++
		time.Sleep(time.Duration(inst.config.PollRate) * time.Second)
		//fmt.Println("\n \n")
		inst.edge28DebugMsg("LOOP COUNT: ", counter)

		nets, err := inst.db.GetNetworksByPlugin(inst.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		if len(nets) == 0 {
			inst.edge28DebugMsg("edge28-polling: NO NETWORKS FOUND")
			time.Sleep(15000 * time.Millisecond)
		}
		for _, net := range nets { // NETWORKS
			inst.edge28DebugMsg("Edge28Polling: net")
			inst.edge28DebugMsg(fmt.Sprintf("%+v\n", net))
			if net.UUID != "" && net.PluginConfId == inst.pluginUUID {
				if boolean.IsFalse(net.Enable) {
					continue
				}
				for _, dev := range net.Devices { // DEVICES
					inst.edge28DebugMsg("Edge28Polling: dev")
					inst.edge28DebugMsg(fmt.Sprintf("%+v\n", dev))
					if err != nil {
						inst.edge28ErrorMsg(fmt.Sprintf("failed to vaildate device %v %s\n", err, dev.CommonIP.Host))
					}
					if boolean.IsFalse(dev.Enable) {
						continue
					}

					rest := edgerest.NewNoAuth(dev.CommonIP.Host, dev.CommonIP.Port)
					getUI, err = rest.GetUIs()
					getDI, err = rest.GetDIs()

					for _, pnt := range dev.Points { // POINTS
						pnt, err := inst.db.GetPoint(pnt.UUID, api.Args{WithPriority: true})
						if err != nil {
							inst.edge28ErrorMsg("cannot find point")
						}

						inst.edge28DebugMsg("Edge28Polling: pnt")
						inst.edge28DebugMsg(fmt.Sprintf("%+v\n", pnt))
						inst.printPointDebugInfo(pnt)
						if boolean.IsFalse(pnt.Enable) {
							inst.edge28DebugMsg("point is disabled.")
							continue
						}
						var rv float64
						var readValStruct interface{}
						var readValType string
						var wv float64

						switch pnt.IoNumber {
						// OUTPUTS
						case pointList.R1, pointList.R2, pointList.DO1, pointList.DO2, pointList.DO3, pointList.DO4, pointList.DO5:
							pnt.PointPriorityArrayMode = model.PriorityArrayToWriteValue
							pnt.WriteValue = limitValueByEdge28Type(pnt.IoType, pnt.WriteValue)
							if pnt.WriteValue != nil {
								writeValue := float.NonNil(pnt.WriteValue)
								wv, err = DigitalToGPIOValue(writeValue, false)
								if err != nil {
									inst.edge28ErrorMsg("invalid input to DigitalToGPIOValue")
									continue
								}
								_, err = inst.processWrite(pnt, wv, rest, uint64(counter), false)
								if err != nil {
									inst.edge28ErrorMsg(err)
									inst.pointUpdateErr(pnt, err)
									continue
								}
							}

						case pointList.UO1, pointList.UO2, pointList.UO3, pointList.UO4, pointList.UO5, pointList.UO6, pointList.UO7:
							pnt.PointPriorityArrayMode = model.PriorityArrayToWriteValue
							pnt.WriteValue = limitValueByEdge28Type(pnt.IoType, pnt.WriteValue)
							if pnt.WriteValue != nil {
								wv, err = GetGPIOValueForUOByType(pnt)
								if err != nil {
									inst.edge28ErrorMsg(err)
									continue
								}
								_, err = inst.processWrite(pnt, wv, rest, uint64(counter), true)
								if err != nil {
									inst.edge28ErrorMsg(err)
									continue
								}
							}

						// INPUTS
						case pointList.DI1, pointList.DI2, pointList.DI3, pointList.DI4, pointList.DI5, pointList.DI6, pointList.DI7:
							pnt.PointPriorityArrayMode = model.ReadOnlyNoPriorityArrayRequired
							if getDI == nil {
								inst.edge28ErrorMsg("error on DI read")
								continue
							}
							readValStruct, readValType, err = structs.GetStructFieldByString(getDI.Val, pnt.IoNumber)
							if err != nil {
								inst.edge28ErrorMsg(err)
								continue
							} else if readValType != "struct" {
								inst.edge28ErrorMsg("IoNumber does not match any points from Edge28")
								continue
							}
							rv = reflect.ValueOf(readValStruct).FieldByName("Val").Float()
							rv, err = GetValueFromGPIOForUIByType(pnt, rv)
							if err != nil {
								inst.edge28ErrorMsg(err)
								continue
							}
							_, err = inst.processRead(pnt, rv, counter)

						case pointList.UI1, pointList.UI2, pointList.UI3, pointList.UI4, pointList.UI5, pointList.UI6, pointList.UI7:
							pnt.PointPriorityArrayMode = model.ReadOnlyNoPriorityArrayRequired
							if getUI == nil {
								inst.edge28ErrorMsg("error on UI read")
								continue
							}
							readValStruct, readValType, err = structs.GetStructFieldByString(getUI.Val, pnt.IoNumber)
							if err != nil {
								inst.edge28ErrorMsg(err)
								continue
							} else if readValType != "struct" {
								inst.edge28ErrorMsg("IoNumber does not match any points from Edge28")
								continue
							}
							rv = reflect.ValueOf(readValStruct).FieldByName("Val").Float()
							rv, err = GetValueFromGPIOForUIByType(pnt, rv)
							if err != nil {
								inst.edge28ErrorMsg(err)
								continue
							}
							_, err = inst.processRead(pnt, rv, counter)
						}
					}
				}
			}
		}
		return false, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	inst.pollingCancel = cancel
	go poll.GoPoll(ctx, f)
	return nil
}

// TODO add COV WriteValueOnceSync and InSync
func (inst *Instance) processWrite(pnt *model.Point, value float64, rest *edgerest.RestClient, pollCount uint64, isUO bool) (float64, error) {
	rsyncWrite := pollCount % 10
	// rsyncWrite is just a way to make sure the outputs on the device are not out of sync
	// and rsync on first poll loop
	writeValue := float.NonNil(pnt.WriteValue)
	var err error
	inst.edge28DebugMsg("processWrite: pnt")
	inst.edge28DebugMsg(fmt.Sprintf("%+v\n", pnt))
	if boolean.IsTrue(pnt.WritePollRequired) || rsyncWrite == 0 || pollCount == 1 {
		if pollCount == 1 {
			inst.edge28DebugMsg(fmt.Sprintf("processWrite() SYNC on first poll wrote IO %s: %v\n", pnt.IoNumber, value))
		}
		if rsyncWrite == 0 {
			inst.edge28DebugMsg(fmt.Sprintf("processWrite() rsyncWrite wrote IO %s: %v\n", pnt.IoNumber, value))
		}
		if isUO {
			inst.edge28DebugMsg(fmt.Sprintf("processWrite() WRITE UO %s as type:%s  value:%v\n", pnt.IoNumber, pnt.IoType, value))
			_, err = rest.WriteUO(pnt.IoNumber, value)
			if err != nil {
				inst.edge28ErrorMsg(fmt.Sprintf("processWrite() failed to write IO %s:  value:%f error:%v\n", pnt.IoNumber, value, err))
				inst.pointUpdateErr(pnt, err)
				return 0, err
			} else {
				_, err = inst.pointUpdate(pnt, writeValue, true, true, true)
			}
		} else {
			inst.edge28DebugMsg(fmt.Sprintf("processWrite() WRITE DO %s as type:%s  value:%v\n", pnt.IoNumber, pnt.IoType, value))
			_, err = rest.WriteDO(pnt.IoNumber, value)
			if err != nil {
				inst.edge28ErrorMsg(fmt.Sprintf("processWrite() failed to write IO %s:  value:%f error:%v\n", pnt.IoNumber, value, err))
				inst.pointUpdateErr(pnt, err)
				return 0, err
			} else {
				_, err = inst.pointUpdate(pnt, writeValue, true, true, true)
			}
		}
		if err != nil {
			inst.edge28ErrorMsg(fmt.Sprintf("processWrite() failed to write IO %s:  value:%f error:%v\n", pnt.IoNumber, value, err))
			err := inst.pointUpdateErr(pnt, err)
			return 0, err
		} else {
			// log.Infof("edge28-polling: wrote IO %s: %v\n", pnt.IoNumber, value)
			return value, err
		}
	} else {
		inst.edge28DebugMsg(fmt.Sprintf("point is in sync %s: %v\n", pnt.IoNumber, value))
		return value, err
	}
}

func (inst *Instance) processRead(pnt *model.Point, readValue float64, pollCount int) (float64, error) {
	covEvent, _ := nmath.Cov(readValue, float.NonNil(pnt.OriginalValue), 0) // TODO: Remove this as it's done in the main point db file
	inst.edge28DebugMsg("processRead: pnt")
	inst.edge28DebugMsg(fmt.Sprintf("%+v\n", pnt))
	if pollCount == 1 || boolean.IsTrue(pnt.ReadPollRequired) {
		_, err = inst.pointUpdate(pnt, readValue, true, true, true)
		if err != nil {
			inst.edge28DebugMsg(fmt.Sprintf("READ UPDATE POINT %s: %v\n", pnt.IoNumber, readValue))
			err := inst.pointUpdateErr(pnt, err)
			return readValue, err
		}
		if boolean.IsTrue(pnt.InSync) {
			inst.edge28DebugMsg(fmt.Sprintf("READ POINT SYNC %s: %v\n", pnt.IoNumber, readValue))
		} else {
			inst.edge28DebugMsg(fmt.Sprintf("READ ON START %s: %v\n", pnt.IoNumber, readValue))
		}
	} else if covEvent {
		_, err = inst.pointUpdate(pnt, readValue, true, true, true)
		if err != nil {
			inst.edge28ErrorMsg(fmt.Sprintf("READ UPDATE POINT %s: %v\n", pnt.IoNumber, readValue))
			err := inst.pointUpdateErr(pnt, err)
			return readValue, err
		} else {
			inst.edge28ErrorMsg(fmt.Sprintf("READ ON START %s: %v\n", pnt.IoNumber, readValue))
		}
	}
	return readValue, nil
}
