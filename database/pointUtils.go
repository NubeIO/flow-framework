package database

import (
	"context"
	"fmt"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	unit "github.com/NubeIO/flow-framework/src/units"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/PaesslerAG/gval"
)

func (d *GormDatabase) pointNameExists(body *model.Point) (nameExist, addressIDExist bool) {
	var arg api.Args
	arg.WithPoints = true
	device, err := d.GetDevice(body.DeviceUUID, arg)
	if err != nil {
		return false, false
	}

	nameExist = false
	addressIDExist = false

	for _, pnt := range device.Points {
		if pnt.Name == body.Name {
			if pnt.UUID == body.UUID {
				nameExist = false
			} else {
				nameExist = true
			}
		}
	}

	//check body with existing devices that a point with the same objectType do not have the same addrID
	for _, pnt := range device.Points {
		if pnt.ObjectType == body.ObjectType {
			if utils.IntIsNil(pnt.AddressID) == utils.IntIsNil(body.AddressID) {
				if pnt.UUID == body.UUID {
					addressIDExist = false
				} else {
					addressIDExist = true
				}
			}
		}
	}
	return nameExist, addressIDExist
}

// PointDeviceByAddressID will query by device_uuid = ? AND object_type = ? AND address_id = ?
func (d *GormDatabase) PointDeviceByAddressID(pointUUID string, body *model.Point) (*model.Point, bool) {
	var pointModel *model.Point
	deviceUUID := body.DeviceUUID
	objType := body.ObjectType
	addressID := body.AddressID
	f := fmt.Sprintf("device_uuid = ? AND object_type = ? AND address_id = ?")
	query := d.DB.Where(f, deviceUUID, objType, addressID, pointUUID).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, false
	}
	return pointModel, true
}

func pointUnits(presentValue *float64, unitFrom, unitTo *string) (value *float64, err error) {
	if presentValue == nil || unitFrom == nil || unitTo == nil || *unitFrom == "" || *unitTo == "" {
		return presentValue, nil
	}
	_, res, err := unit.Process(*presentValue, *unitFrom, *unitTo)
	if err != nil {
		return utils.NewFloat64(0), err
	}
	return utils.NewFloat64(res.AsFloat()), err
}

func pointRange(presentValue, limitMin, limitMax *float64) (value *float64) {
	if !utils.FloatIsNilCheck(presentValue) && !utils.FloatIsNilCheck(limitMin) && !utils.FloatIsNilCheck(limitMax) {
		if *limitMin == 0 && *limitMax == 0 {
			return presentValue
		}
		out := utils.LimitToRange(*presentValue, *limitMin, *limitMax)
		return &out
	}
	return presentValue
}

func pointScale(presentValue, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax *float64) (value *float64) {
	if !utils.FloatIsNilCheck(presentValue) && !utils.FloatIsNilCheck(scaleInMin) && !utils.FloatIsNilCheck(scaleInMax) && !utils.FloatIsNilCheck(scaleOutMin) && !utils.FloatIsNilCheck(scaleOutMax) {
		if *scaleInMin == 0 && *scaleInMax == 0 && *scaleOutMin == 0 && *scaleOutMax == 0 {
			return presentValue
		}
		out := utils.Scale(*presentValue, *scaleInMin, *scaleInMax, *scaleOutMin, *scaleOutMax)
		return &out
	}
	return presentValue
}

func pointEval(val *float64, evalString string) (value *float64, err error) {
	if evalString != "" {
		eval, err := gval.Full().NewEvaluable(evalString)
		if err != nil || val == nil {
			return nil, err
		}
		v, err := eval.EvalFloat64(context.Background(), map[string]interface{}{"x": *val})
		if err != nil {
			return nil, err
		}
		_v := v
		return &_v, nil
	}
	return val, nil
}
