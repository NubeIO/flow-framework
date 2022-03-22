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

//TODO add COV WriteValueOnceSync and InSync
func (i *Instance) processWrite(pnt *model.Point, value float64, rest *edgerest.RestClient, pollCount uint64, isUO bool) (float64, error) {
	rsyncWrite := pollCount % 10
	//rsyncWrite is just a way to make sure the outputs on the device are not out of sync
	// and rsync on first poll loop
	var err error
	if !utils.BoolIsNil(pnt.InSync) || rsyncWrite == 0 || pollCount == 1 {
		if pollCount == 1 {
			log.Infof("edge28-polling: SYNC on first poll wrote IO %s: %v\n", pnt.IoNumber, value)
		}
		if rsyncWrite == 0 {
			log.Infof("edge28-polling: rsyncWrite wrote IO %s: %v\n", pnt.IoNumber, value)
		}
		if isUO {
			_, err = rest.WriteUO(pnt.IoNumber, value)
			_, err = i.pointUpdateValue(pnt.UUID, value)
		} else {
			_, err = rest.WriteDO(pnt.IoNumber, value)
			_, err = i.pointUpdateValue(pnt.UUID, value)
		}
		if err != nil {
			log.Errorf("edge28-polling: failed to write IO %s:  value:%f error:%v\n", pnt.IoNumber, value, err)
			_, err := i.pointUpdateErr(pnt.UUID, err)
			return 0, err
		} else {
			log.Infof("edge28-polling: wrote IO %s: %v\n", pnt.IoNumber, value)
			return value, err
		}
	} else {
		log.Infof("edge28-polling: point is in sync %s: %v\n", pnt.IoNumber, value)
		return value, err
	}

}

func (i *Instance) processRead(pnt *model.Point, readValue float64, pollCount float64) (float64, error) {
	covEvent, _ := utils.COV(readValue, utils.Float64IsNil(pnt.PresentValue), 0) //Remove this as it's done in the main point db file
	if pollCount == 1 || !utils.BoolIsNil(pnt.InSync) {
		_, err := i.pointUpdateValue(pnt.UUID, readValue)
		if err != nil {
			log.Errorf("edge28-polling: READ UPDATE POINT %s: %v\n", pnt.IoNumber, readValue)
			_, err := i.pointUpdateErr(pnt.UUID, err)
			return readValue, err
		}
		if utils.BoolIsNil(pnt.InSync) {
			log.Infof("edge28-polling: READ POINT SYNC %s: %v\n", pnt.IoNumber, readValue)
		} else {
			log.Infof("edge28-polling: READ ON START %s: %v\n", pnt.IoNumber, readValue)
		}
	} else if covEvent {
		_, err := i.pointUpdateValue(pnt.UUID, readValue)
		if err != nil {
			log.Errorf("edge28-polling READ UPDATE POINT %s: %v\n", pnt.IoNumber, readValue)
			_, err := i.pointUpdateErr(pnt.UUID, err)
			return readValue, err
		} else {
			log.Infof("edge28-polling: READ ON START %s: %v\n", pnt.IoNumber, readValue)
		}
	}
	return readValue, nil
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
			log.Info("edge28-polling: NO NETWORKS FOUND")
		}
		//log.Error(bugs.DebugPrint(name+pollName, i.polling, err, "hey", "whats", "up"))
		for _, net := range nets { //NETWORKS
			if net.UUID != "" && net.PluginConfId == i.pluginUUID {
				log.Infof("edge28-polling: LOOP COUNT: %v\n", counter)
				counter++
				for _, dev := range net.Devices { //DEVICES
					if err != nil {
						log.Errorf("edge28-polling: failed to vaildate device %v %s\n", err, dev.CommonIP.Host)
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
						switch pnt.IoNumber {
						//OUTPUTS
						case pointList.R1, pointList.R2, pointList.DO1, pointList.DO2, pointList.DO3, pointList.DO4, pointList.DO5:
							//_, err := i.db.UpdatePointValue(pnt.UUID, pnt, false)  //TODO: This call sets the fallback value, but it ends up being called too often and overrides value changes from API calls
							if pnt.PresentValue != nil {
								wv, err = DigitalToGPIOValue(*(pnt.PresentValue), false)
								if err != nil {
									log.Errorf("edge28-polling: invalid input to  DigitalToGPIOValue")
									continue
								}
								_, err = i.processWrite(pnt, wv, rest, uint64(counter), false)
							}
						case pointList.UO1, pointList.UO2, pointList.UO3, pointList.UO4, pointList.UO5, pointList.UO6, pointList.UO7:
							//_, err := i.db.UpdatePointValue(pnt.UUID, pnt, false) //TODO: This call sets the fallback value, but it ends up being called too often and overrides value changes from API calls
							if pnt.PresentValue != nil {
								wv, err = GetGPIOValueForUOByType(pnt)
								if err != nil {
									log.Error(err)
									continue
								}
								_, err = i.processWrite(pnt, wv, rest, uint64(counter), true)
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
							readValStruct, readValType, err = utils.GetStructFieldByString(getDI.Val, pnt.IoNumber)
							if err != nil {
								log.Error("edge28-polling: ", err)
								continue
							} else if readValType != "struct" {
								log.Error("edge28-polling: IoNumber does not match any points from Edge28")
								continue
							}
							rv = reflect.ValueOf(readValStruct).FieldByName("Val").Float()
							rv, err = GetValueFromGPIOForUIByType(pnt, rv)
							if err != nil {
								log.Error("edge28-polling: ", err)
								continue
							}
							_, err = i.processRead(pnt, rv, counter)

						case pointList.UI1, pointList.UI2, pointList.UI3, pointList.UI4, pointList.UI5, pointList.UI6, pointList.UI7:
							if getUI == nil {
								continue
							}
							readValStruct, readValType, err = utils.GetStructFieldByString(getUI.Val, pnt.IoNumber)
							if err != nil {
								log.Error(err)
								continue
							} else if readValType != "struct" {
								log.Error("edge28-polling: IoNumber does not match any points from Edge28")
								continue
							}
							rv = reflect.ValueOf(readValStruct).FieldByName("Val").Float()
							rv, err = GetValueFromGPIOForUIByType(pnt, rv)
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
		err = errors.New(fmt.Sprintf("edge-28: skipping %v, IoType %v not recognized.", point.IoNumber, point.IoType))
		return 0, err
	}
	//wv = *(point.PresentValue)   //TODO: use PresentValue instead of Priority 16 value
	if numbers.Float64PointerIsNil(point.PresentValue) {
		return 0, errors.New("edge28-polling: no value to write")
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
		err = errors.New("edge28-polling: UO IoType is not a recognized type")
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
		err = errors.New(fmt.Sprintf("edge28-polling: skipping %v, IoType %v not recognized.", point.IoNumber, point.IoType))
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
		err = errors.New("edge28-polling: UI IoType is not a recognized type")
		return 0, err
	}
	return result, nil
}

//DigitalToGPIOValue converts true/false values (all basic types allowed) to BBB GPIO 0/1 ON/OFF (FOR DOs/Relays) and to 100/0 (FOR UOs).  Note that the GPIO value for digital points is inverted.
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
		err = errors.New("edge28-polling: input is not a recognized type")
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
