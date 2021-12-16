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

func (d *GormDatabase) pointNameExists(pnt *model.Point, body *model.Point) bool {
	var arg api.Args
	arg.WithPoints = true
	device, err := d.GetDevice(pnt.DeviceUUID, arg)
	if err != nil {
		return false
	}
	if pnt.UUID == "" && body.UUID == "" {
		return false
	}
	for _, p := range device.Points {
		if p.Name == body.Name {
			if p.UUID == pnt.UUID {
				return false
			} else {
				return true
			}
		}
	}
	return false
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

func pointUnits(pointModel *model.Point) (value float64, ok bool, err error) {
	if pointModel.Unit != "" {
		_, res, err := unit.Process(*pointModel.PresentValue, pointModel.Unit, pointModel.UnitTo)
		if err != nil {
			return 0, false, err
		}
		return res.AsFloat(), true, err
	} else {
		return 0, false, nil
	}
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

func pointEval(presentValue, originalValue *float64, evalMode, evalString string) (value *float64, err error) {

	var val *float64
	if model.EvalMode(evalMode) == model.EvalModeCalcAfterScale || model.EvalMode(evalMode) == model.EvalModeEnable {
		val = presentValue
	} else if model.EvalMode(evalMode) == model.EvalModeCalcOnOriginalValue {
		val = originalValue
	} else {
		val = presentValue
	}
	exp := evalString
	if evalString != "" && model.EvalMode(evalMode) != model.EvalModeDisabled {
		eval, err := gval.Full().NewEvaluable(exp)
		if err != nil && val != nil {
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
