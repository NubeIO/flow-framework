package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	edgerest "github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/restclient"
	"github.com/NubeIO/flow-framework/src/poller"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/edge28"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/numbers"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/thermistor"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"time"
)

const defaultInterval = 2500 * time.Millisecond //default polling is 2.5 sec

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

//TODO add COV WriteValueOnceSync and InSync
func (i *Instance) processWrite(pnt *model.Point, value float64, rest *edgerest.RestClient, pollCount float64, isUO bool) (float64, error) {
	//!utils.BoolIsNil(pnt.WriteValueOnceSync)
	var err error
	if isUO {
		_, err = rest.WriteUO(pnt.IoID, value)
	} else {
		_, err = rest.WriteDO(pnt.IoID, value)
	}
	if err != nil {
		log.Errorf("edge-28: failed to write IO %s:  value:%f error:%v\n", pnt.IoID, value, err)
		return 0, err
	} else {
		log.Infof("edge-28: wrote IO %s: %v\n", pnt.IoID, value)
		return value, err
	}
}

func (i *Instance) processRead(pnt *model.Point, value float64, pollCount float64) (float64, error) {
	cov := utils.Float64IsNil(pnt.COV) //TODO add in point scaling to get COV to work correct (as in scale temp or 0-10)
	covEvent, _ := utils.COV(value, utils.Float64IsNil(pnt.PresentValue), cov)
	if pollCount == 1 || !utils.BoolIsNil(pnt.InSync) {
		pnt.InSync = utils.NewTrue()
		pnt.Priority.P16 = utils.NewFloat64(value)
		_, err := i.db.UpdatePointValue(pnt.UUID, pnt, false)
		if err != nil {
			log.Errorf("edge-28: READ UPDATE POINT %s: %v\n", pnt.IoID, value)
			return value, err
		}
		if utils.BoolIsNil(pnt.InSync) {
			log.Infof("edge-28: READ POINT SYNC %s: %v\n", pnt.IoID, value)
		} else {
			log.Infof("edge-28: READ ON START %s: %v\n", pnt.IoID, value)
		}
	} else if covEvent {
		pnt.InSync = utils.NewTrue()
		fmt.Println("processRead()  value:", value)
		fmt.Println("pnt.Priority.P16 - 1", *(pnt.Priority.P16))
		pnt.Priority.P16 = utils.NewFloat64(value)
		fmt.Printf("%+v\n", pnt)
		fmt.Printf("%+v\n", pnt.Priority.P16)
		fmt.Println("pnt.Priority.P16 - 2", *(pnt.Priority.P16))
		_, err := i.db.UpdatePointValue(pnt.UUID, pnt, true)
		if err != nil {
			log.Errorf("edge-28: READ UPDATE POINT %s: %v\n", pnt.IoID, value)
			return value, err
		} else {
			log.Infof("edge-28: READ ON START %s: %v\n", pnt.IoID, value)
		}
	}
	return value, nil
}

func (i *Instance) polling(p polling) error {
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
		nets, err := i.db.GetNetworksByPlugin(i.pluginUUID, arg)
		if err != nil {
			return false, err
		}
		if len(nets) == 0 {
			time.Sleep(15000 * time.Millisecond)
			log.Info("edge-28: NO NETWORKS FOUND")
		}

		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				log.Infof("edge-28: LOOP COUNT: %v\n", counter)
				counter++

				for _, dev := range net.Devices { //DEVICES
					if err != nil {
						log.Errorf("edge-28: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
					}
					rest := edgerest.NewNoAuth(dev.CommonIP.Host, dev.CommonIP.Port)
					getUI, err = rest.GetUIs()
					getDI, err = rest.GetDIs()
					dNet := p.delayNetworks
					time.Sleep(dNet)

					for _, pnt := range dev.Points { //POINTS
						var rv float64
						var readValStruct interface{}
						var readValType string
						var wv float64
						fmt.Println(pnt.IoID)
						switch pnt.IoID {
						//OUTPUTS
						case pointList.R1, pointList.R2, pointList.DO1, pointList.DO2, pointList.DO3, pointList.DO4, pointList.DO5:
							//_, err := i.db.UpdatePointValue(pnt.UUID, pnt, false)  //TODO: This call sets the fallback value, but it ends up being called too often and overrides value changes from API calls
							if pnt.PresentValue != nil {
								wv, err = DigitalToGPIOValue(*(pnt.PresentValue), false)
								if err != nil {
									log.Errorf("edge-28: invalid input to  DigitalToGPIOValue")
									continue
								}
								fmt.Println("Relay/DO to processWrite():", wv)
								_, err = i.processWrite(pnt, wv, rest, counter, false)
							}

						case pointList.UO1, pointList.UO2, pointList.UO3, pointList.UO4, pointList.UO5, pointList.UO6, pointList.UO7:
							//fmt.Println(*(pnt))
							//fmt.Printf("%+v\n", *(pnt))
							//_, err := i.db.UpdatePointValue(pnt.UUID, pnt, false) //TODO: This call sets the fallback value, but it ends up being called too often and overrides value changes from API calls
							if pnt.PresentValue != nil {
								wv, err = GetGPIOValueForUOByType(pnt)
								if err != nil {
									log.Error(err)
									continue
								}
								_, err = i.processWrite(pnt, wv, rest, counter, true)
								if err != nil {
									log.Error(err)
									continue
								}
							}

						//INPUTS
						case pointList.DI1, pointList.DI2, pointList.DI3, pointList.DI4, pointList.DI5, pointList.DI6, pointList.DI7:
							if getDI == nil {
								continue
							}
							readValStruct, readValType, err = utils.GetStructFieldByString(getDI.Val, pnt.IoID)
							if err != nil {
								log.Error(err)
								continue
							} else if readValType != "struct" {
								log.Error("edge-28: IoID does not match any points from Edge28")
								continue
							}
							rv = reflect.ValueOf(readValStruct).FieldByName("Val").Float()
							rv, err = GetValueFromGPIOForUIByType(pnt, rv)
							if err != nil {
								log.Error(err)
								continue
							}
							_, err = i.processRead(pnt, rv, counter)

						case pointList.UI1, pointList.UI2, pointList.UI3, pointList.UI4, pointList.UI5, pointList.UI6, pointList.UI7:
							if getUI == nil {
								continue
							}
							fmt.Println("POINT")
							fmt.Printf("%+v\n", *(pnt))
							readValStruct, readValType, err = utils.GetStructFieldByString(getUI.Val, pnt.IoID)
							fmt.Println("readValStruct", readValStruct)
							fmt.Println("readValType", readValType)
							//fmt.Printf("%+v\n", *(pnt))
							if err != nil {
								log.Error(err)
								continue
							} else if readValType != "struct" {
								log.Error("edge-28: IoID does not match any points from Edge28")
								continue
							}
							rv = reflect.ValueOf(readValStruct).FieldByName("Val").Float()
							fmt.Println("rv1", rv)
							rv, err = GetValueFromGPIOForUIByType(pnt, rv)
							fmt.Println("rv2", rv)
							if err != nil {
								log.Error(err)
								continue
							}
							_, err = i.processRead(pnt, rv, counter)
						}
					}
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
		return err
	}
	return nil
}

// GetGPIOValueForUOByType converts the point value to the correct edge28 UO GPIO value based on the IoType
func GetGPIOValueForUOByType(point *model.Point) (float64, error) {
	var err error
	var result float64
	if !utils.ExistsInStrut(UOTypes, point.IoType) {
		err = errors.New(fmt.Sprintf("edge-28: skipping %v, IoType %v not recognized.", point.IoID, point.IoType))
		return 0, err
	}
	//fmt.Println("point")
	//fmt.Printf("%+v\n", point)
	//fmt.Println("point.Priority")
	//fmt.Printf("%+v\n", point.Priority)
	//fmt.Println("point.Priority.P16")
	//fmt.Printf("%+v\n", point.Priority.P16)
	//wv = *(point.PresentValue)   //TODO: use PresentValue instead of Priority 16 value
	if numbers.Float64PointerIsNil(point.PresentValue) {
		return 0, errors.New("no value to write.")
	} else {
		result = *(point.PresentValue)
	}

	switch point.IoType {
	case UOTypes.DIGITAL:
		result, err = DigitalToGPIOValue(result, true)
	case UOTypes.PERCENT:
		result = edge28.PercentToGPIOValue(result)
	case UOTypes.VOLTSDC:
		result = edge28.VoltageToGPIOValue(result)
	default:
		err = errors.New("UO IoType is not a recognized type")
	}
	if err != nil {
		return 0, err
	} else {
		return result, nil
	}
}

// GetValueFromGPIOForUIByType converts the GPIO value to the scaled UI value based on the IoType
func GetValueFromGPIOForUIByType(point *model.Point, value float64) (float64, error) {
	var err error
	var result float64

	if !utils.ExistsInStrut(UITypes, point.IoType) {
		err = errors.New(fmt.Sprintf("edge-28: skipping %v, IoType %v not recognized.", point.IoID, point.IoType))
		return 0, err
	}
	switch point.IoType {
	case UITypes.RAW:
		result = value
	case UITypes.DIGITAL:
		result = edge28.GPIOValueToDigital(value)
	case UITypes.PERCENT:
		result = edge28.GPIOValueToPercent(value)
	case UITypes.VOLTSDC:
		result = edge28.GPIOValueToVoltage(value)
	case UITypes.MILLIAMPS:
		result = edge28.ScaleGPIOValueTo420ma(value)
	case UITypes.RESISTANCE:
		result = edge28.ScaleGPIOValueToResistance(value)
	case UITypes.THERMISTOR10KT2:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.T210K)
	case UITypes.THERMISTOR10KT3:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.T310K)
	case UITypes.THERMISTOR20KT1:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.T120K)
	case UITypes.THERMISTORPT100:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.PT100)
	case UITypes.THERMISTORPT1000:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.PT1000)
	default:
		err = errors.New("UI IoType is not a recognized type")
		return 0, err
	}
	fmt.Println("result", result)
	return result, nil
}

//TODO: update this function in Edge28 helpers
//DigitalToRelayGPIO converts true/false values (all basic types allowed) to BBB GPIO 0/1 ON/OFF (FOR DOs/Relays) and to 100/0 (FOR UOs).  Note that the GPIO value for digital points is inverted.
func DigitalToGPIOValue(input interface{}, isUO bool) (float64, error) {
	var inputAsBool bool
	var err error = nil
	switch input.(type) {
	case string:
		inputAsBool, err = strconv.ParseBool(reflect.ValueOf(input).String())
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		inputAsBool = reflect.ValueOf(input).Int() != 0
	case float32, float64:
		inputAsBool = reflect.ValueOf(input).Float() != float64(0)
	case bool:
		inputAsBool = reflect.ValueOf(input).Bool()
	default:
		err = errors.New("input is not a recognized type")
	}
	if err != nil {
		return 0, err
	} else if inputAsBool {
		if isUO {
			return 0, nil // 0 is the 12vdc/ON GPIO value for UOs
		} else {
			return 1, nil // 0 is the 12vdc/ON GPIO value for DOs/Relays
		}
	} else {
		if isUO {
			return 100, nil // 100 is the 0vdc/OF GPIO value for UOs
		} else {
			return 0, nil // 1 is the 0vdc/OFF GPIO value for DOs/Relays
		}
	}
}
