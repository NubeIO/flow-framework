package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	log "github.com/sirupsen/logrus"
)

func (d *GormDatabase) GetPoints(args api.Args) ([]*model.Point, error) {
	var pointsModel []*model.Point
	query := d.buildPointQuery(args)
	if err := query.Find(&pointsModel).Error; err != nil {
		return nil, err
	}
	return pointsModel, nil
}

func (d *GormDatabase) GetPoint(uuid string, args api.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.Where("uuid = ? ", uuid).First(&pointModel).Error; err != nil {
		return nil, err
	}
	return pointModel, nil
}

func (d *GormDatabase) CreatePoint(body *model.Point, fromPlugin bool) (*model.Point, error) {
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Point)
	deviceUUID := body.DeviceUUID
	body.Name = nameIsNil(body.Name)
	existingAddrID := false
	existingName, _ := d.pointNameExists(body)
	if body.AddressID != nil {
		_, existingAddrID = d.pointNameExists(body)
	}
	if existingName {
		eMsg := fmt.Sprintf("a point with existing name: %s exists", body.Name)
		return nil, errors.New(eMsg)
	}
	if existingAddrID {
		eMsg := fmt.Sprintf("a point with existing AddressID: %d exists", utils.IntIsNil(body.AddressID))
		return nil, errors.New(eMsg)
	}
	if body.Decimal == nil {
		body.Decimal = nils.NewUint32(2)
	}

	obj, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	body.ObjectType = string(obj)
	if body.Description == "" {
		body.Description = "na"
	}
	body.ThingClass = model.ThingClass.Point
	body.CommonEnable.Enable = utils.NewTrue()
	body.InSync = utils.NewFalse()
	body.WriteValueOnceSync = utils.NewFalse()
	if body.Priority == nil {
		body.Priority = &model.Priority{}
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	plug, err := d.GetPluginIDFromDevice(deviceUUID)
	if err != nil {
		return nil, errors.New("ERROR failed to get plugin uuid")
	}
	if !fromPlugin {
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, plug.PluginConfId, body.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, body)
		if err != nil {
			return nil, errors.New("ERROR on device eventbus")
		}
	}
	return body, err
}

func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	existingName, existingAddrID := d.pointNameExists(body)
	if existingName {
		eMsg := fmt.Sprintf("a point with existing name: %s exists", body.Name)
		return nil, errors.New(eMsg)
	}

	if !utils.IntNilCheck(body.AddressID) {
		if existingAddrID {
			eMsg := fmt.Sprintf("a point with existing AddressID: %d exists", utils.IntIsNil(body.AddressID))
			return nil, errors.New(eMsg)
		}
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&pointModel, body.Tags); err != nil {
			return nil, err
		}
	}
	//example modbus: if user changes the data type then do a new read of the point on the modbus network
	if !fromPlugin {
		pointModel.InSync = utils.NewFalse()
	}
	body.WriteValueOnceSync = utils.NewFalse()
	query = d.DB.Model(&pointModel).Updates(&body)
	// Don't update point value if priority array on body is nil
	if body.Priority == nil {
		return pointModel, nil
	} else {
		pointModel.Priority = body.Priority
	}

	pnt, err := d.UpdatePointValue(pointModel, fromPlugin)
	if err != nil {
		return nil, err
	}
	return pnt, nil
}

func (d *GormDatabase) PointWrite(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Priority == nil {
		return nil, errors.New("no priority value is been sent")
	} else {
		pointModel.Priority = body.Priority
	}
	pointModel.InSync = utils.NewFalse()
	point, err := d.UpdatePointValue(pointModel, fromPlugin)
	return point, err
}

func (d *GormDatabase) UpdatePointValue(pointModel *model.Point, fromPlugin bool) (*model.Point, error) {
	pointModel, presentValue := d.updatePriority(pointModel)

	ov := utils.Float64IsNil(presentValue)
	pointModel.OriginalValue = &ov

	presentValue = pointScale(presentValue, pointModel.ScaleInMin, pointModel.ScaleInMax, pointModel.ScaleOutMin, pointModel.ScaleOutMax)
	presentValue = pointRange(presentValue, pointModel.LimitMin, pointModel.LimitMax)
	eval, err := pointEval(presentValue, pointModel.MathOnPresentValue)
	if err != nil {
		log.Errorln("point.db UpdatePointValue() error on run point MathOnPresentValue error:", err)
		return nil, err
	} else {
		presentValue = eval
	}
	val, err := pointUnits(presentValue, pointModel.Unit, pointModel.UnitTo)
	if err != nil {
		log.Errorf("ERROR on point invalid point unit")
		return nil, err
	}
	presentValue = val
	//example for wires and modbus: if a new value is written from  wires then set this to false so the modbus knows on the next poll to write a new value to the modbus point
	if !fromPlugin {
		pointModel.InSync = utils.NewFalse()
	}
	if !utils.Unit32NilCheck(pointModel.Decimal) && presentValue != nil {
		val := utils.RoundTo(*presentValue, *pointModel.Decimal)
		presentValue = &val
	}

	isChange := !utils.CompareFloatPtr(pointModel.PresentValue, presentValue)
	pointModel.PresentValue = presentValue
	if presentValue == nil {
		// nil is ignored on GORM, so we are pushing forcefully because isChange comparison will fail on `null` write
		d.DB.Model(&pointModel).Update("present_value", nil)
		d.DB.Model(&model.Writer{}).
			Where("writer_thing_uuid = ?", pointModel.UUID).
			Update("present_value", nil)
	}
	if isChange == true {
		_ = d.DB.Model(&pointModel).Updates(&pointModel)
		err = d.ProducersPointWrite(pointModel)
		if err != nil {
			return nil, err
		}
		d.DB.Model(&model.Writer{}).
			Where("writer_thing_uuid = ?", pointModel.UUID).
			Update("present_value", pointModel.PresentValue)
	}

	if !fromPlugin { // stop looping
		plug, err := d.GetPluginIDFromDevice(pointModel.DeviceUUID)
		if err != nil {
			return nil, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, plug.PluginConfId, pointModel.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, pointModel)
		if err != nil {
			return nil, errors.New("ERROR on device eventbus")
		}
	}
	return pointModel, nil
}

func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	var pointModel *model.Point
	point, err := d.GetPoint(uuid, api.Args{})
	if err != nil {
		return false, errors.New("point not exist")
	}
	query := d.DB.Where("uuid = ? ", uuid).Delete(&pointModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		plug, err := d.GetPluginIDFromDevice(point.DeviceUUID)
		if err != nil {
			return false, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsDeleted, plug.PluginConfId, point.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, point)
		if err != nil {
			return false, errors.New("ERROR on device eventbus")
		}
		return true, nil
	}
}

func (d *GormDatabase) DropPoints() (bool, error) {
	var pointModel *model.Point
	query := d.DB.Where("1 = 1").Delete(&pointModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
