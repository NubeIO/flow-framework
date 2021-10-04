package database

import (
	"context"
	"github.com/NubeDev/flow-framework/model"
	unit "github.com/NubeDev/flow-framework/src/units"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/PaesslerAG/gval"
)

func pointUnits(pointModel *model.Point) (value float64, displayValue string, ok bool, err error) {
	if pointModel.Unit != "" {
		_, res, err := unit.Process(*pointModel.PresentValue, pointModel.Unit, pointModel.UnitTo)
		if err != nil {
			return 0, "", false, err
		}
		return res.AsFloat(), res.String(), true, err
	} else {
		return 0, "", false, nil
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

func pointScale(presentValue *float64, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax *float64) (value *float64) {
	if !utils.FloatIsNilCheck(presentValue) && !utils.FloatIsNilCheck(scaleInMin) && !utils.FloatIsNilCheck(scaleInMax) && !utils.FloatIsNilCheck(scaleOutMin) && !utils.FloatIsNilCheck(scaleOutMax) {
		if *scaleInMin == 0 && *scaleInMax == 0 && *scaleOutMin == 0 && *scaleOutMax == 0 {
			return presentValue
		}
		out := utils.Scale(*presentValue, *scaleInMin, *scaleInMax, *scaleOutMin, *scaleOutMax)
		return &out
	}
	return presentValue
}

func pointEval(presentValue, valueOriginal *float64, evalMode, evalString string) (value *float64, err error) {
	var val *float64
	if evalMode == model.EvalMode.CalcAfterScale || evalMode == model.EvalMode.Enable {
		val = presentValue
	} else if evalMode == model.EvalMode.CalcOnValueOriginal {
		val = valueOriginal
	} else {
		val = presentValue
	}
	exp := evalString
	if evalString != "" && evalMode != model.EvalMode.Disabled {
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
