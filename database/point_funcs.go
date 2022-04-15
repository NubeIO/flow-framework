package database

import (
	"errors"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (d *GormDatabase) GetPointByName(networkName, deviceName, pointName string) (*model.Point, error) {
	var pointModel *model.Point
	net, err := d.GetNetworkByName(networkName, api.Args{WithDevices: true, WithPoints: true})
	if err != nil {
		return nil, errors.New("failed to find a network with that name")
	}
	deviceExist := false
	pointExist := false
	for _, device := range net.Devices {
		if device.Name == deviceName {
			deviceExist = true
			for _, p := range device.Points {
				if p.Name == pointName {
					pointExist = true
					pointModel = p
					break
				}
			}
		}
	}
	if !deviceExist {
		return nil, errors.New("failed to find a device with that name")
	}
	if !pointExist {
		return nil, errors.New("found device but failed to find a point with that name")
	}
	return pointModel, nil
}

func (d *GormDatabase) PointWriteByName(networkName, deviceName, pointName string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	point, err := d.GetPointByName(networkName, deviceName, pointName)
	if err != nil {
		return nil, err
	}
	write, err := d.PointWrite(point.UUID, body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return write, nil
}

func (d *GormDatabase) GetOnePointByArgs(args api.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.First(&pointModel).Error; err != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

func (d *GormDatabase) updatePriority(pointModel *model.Point) (*model.Point, *float64) {
	var presentValue *float64
	presentValueFromPriority := pointModel.PointPriorityArrayMode != model.ReadOnlyNoPriorityArrayRequired && pointModel.PointPriorityArrayMode != model.PriorityArrayToWriteValue
	log.Infof("updatePriority() presentValueFromPriority = %t", presentValueFromPriority)
	log.Infof("updatePriority() pointModel = %+v", pointModel)
	// THESE VALUES ARE NOT REQUIRED FOR model.ReadOnlyNoPriorityArrayRequired
	if pointModel.PointPriorityArrayMode == model.ReadOnlyNoPriorityArrayRequired {
		pointModel.CurrentPriority = nil
		pointModel.WriteValue = nil
		pointModel.WriteValueOriginal = nil
		pointModel.Priority.P1 = nil
		pointModel.Priority.P2 = nil
		pointModel.Priority.P3 = nil
		pointModel.Priority.P4 = nil
		pointModel.Priority.P5 = nil
		pointModel.Priority.P6 = nil
		pointModel.Priority.P7 = nil
		pointModel.Priority.P8 = nil
		pointModel.Priority.P9 = nil
		pointModel.Priority.P10 = nil
		pointModel.Priority.P11 = nil
		pointModel.Priority.P12 = nil
		pointModel.Priority.P13 = nil
		pointModel.Priority.P14 = nil
		pointModel.Priority.P15 = nil
		pointModel.Priority.P16 = nil
	}
	if pointModel.Priority != nil {
		priorityMap, highestValue, currentPriority, isPriorityExist := d.parsePriority(pointModel.Priority, pointModel)
		log.Infof("updatePriority() highestValue = %f", highestValue)
		if isPriorityExist {
			pointModel.CurrentPriority = &currentPriority
			pointModel.WriteValueOriginal = utils.NewFloat64(highestValue)
			writeValue, err := pointEval(&highestValue, pointModel.MathOnWriteValue)
			if err != nil {
				log.Errorln("point.db parsePriority() error on run point MathOnWriteValue error:", err)
			} else {
				pointModel.WriteValue = writeValue
			}
			if presentValueFromPriority { // only update presentValue if required by PointPriorityArrayMode
				presentValue = &highestValue
			}
		} else if !utils.FloatIsNilCheck(pointModel.Fallback) {
			pointModel.Priority.P16 = utils.NewFloat64(*pointModel.Fallback)
			pointModel.CurrentPriority = utils.NewInt(16)
			pointModel.WriteValueOriginal = utils.NewFloat64(*pointModel.Priority.P16)
			writeValue, err := pointEval(pointModel.Priority.P16, pointModel.MathOnWriteValue)
			if err != nil {
				log.Errorln("point.db parsePriority() error on run point MathOnWriteValue error:", err)
			} else {
				pointModel.WriteValue = writeValue
			}
			if presentValueFromPriority {
				presentValue = utils.NewFloat64(*pointModel.Fallback) // only update presentValue if required by PointPriorityArrayMode
			}
		}
		//writeValue := utils.Float64IsNil(pointModel.WriteValue)
		d.DB.Model(&model.Point{}).Where("uuid = ?", pointModel.UUID).Update("write_value", pointModel.WriteValue) // TODO: is this neccesary here, should be included in the db call in UpdatePointValue()?
		d.DB.Model(&model.Priority{}).Where("point_uuid = ?", pointModel.UUID).Updates(&priorityMap)
	}
	if !presentValueFromPriority { // if presentValue doesn't come from priority array, update with the existing value (before math exp is applied)
		presentValue = pointModel.OriginalValue
	}
	return pointModel, presentValue
}

func (d *GormDatabase) parsePriority(priority *model.Priority, pointModel *model.Point) (map[string]interface{}, float64, int, bool) {
	priorityMap := map[string]interface{}{}
	priorityValue := reflect.ValueOf(*priority)
	typeOfPriority := priorityValue.Type()
	highestValue := 0.0
	currentPriority := 0
	isPriorityExist := false
	for i := 0; i < priorityValue.NumField(); i++ {
		if priorityValue.Field(i).Type().Kind().String() == "ptr" {
			val := priorityValue.Field(i).Interface().(*float64)
			if val == nil {
				priorityMap[typeOfPriority.Field(i).Name] = nil
			} else {
				if !isPriorityExist {
					currentPriority = i
					highestValue = *val
					writeValue, err := pointEval(val, pointModel.MathOnWriteValue)
					if err != nil {
						log.Errorln("point.db parsePriority() error on run point MathOnWriteValue error:", err)
						//return nil, 0, 0, false
					}
					pointModel.WriteValue = writeValue
					pointModel.WriteValueOriginal = val
				}
				priorityMap[typeOfPriority.Field(i).Name] = *val
				isPriorityExist = true
			}
		}
	}
	return priorityMap, highestValue, currentPriority, isPriorityExist
}
