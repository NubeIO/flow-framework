package database

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
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

// GetPointsByNetworkUUID get all points by a networkUUID, will return all the points under a network
func (d *GormDatabase) GetPointsByNetworkUUID(networkUUID string) (*utils.Array, error) {
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	network, err := d.GetNetwork(networkUUID, arg)
	if err != nil {
		return nil, err
	}
	p := utils.NewArray()
	for _, dev := range network.Devices {
		for _, pnt := range dev.Points {
			p.Add(pnt)
		}
	}
	return p, nil
}

// GetPointsByNetworkPluginName get all points by a network plugin name, will return all the points under a network
func (d *GormDatabase) GetPointsByNetworkPluginName(name string) (*utils.Array, error) {
	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	network, err := d.GetNetworkByPluginName(name, arg)
	if err != nil {
		return nil, err
	}
	points, err := d.GetPointsByNetworkUUID(network.UUID)
	if err != nil {
		return nil, err
	}
	return points, nil
}

//GetPointByName get point by name
func (d *GormDatabase) GetPointByName(networkName, deviceName, pointName string) (*model.Point, error) {
	var args api.Args
	args.WithDevices = true
	args.WithPoints = true
	var pointModel *model.Point
	net, err := d.GetNetworkByName(networkName, args)
	if net.UUID == "" || err != nil {
		return nil, errors.New("failed to find a network with that name")
	}
	foundDev := false
	foundPnt := false
	for _, dev := range net.Devices {
		if dev.Name == deviceName {
			foundDev = true
			for _, pnt := range dev.Points {
				if pnt.Name == pointName {
					foundPnt = true
					pointModel = pnt
				}
			}
		}
	}
	if !foundDev {
		return nil, errors.New("failed to find a device with that name")
	}
	if !foundPnt {
		return nil, errors.New("found device but failed to find a point with that name")
	}
	return pointModel, nil
}

//PointWriteByName get point by name and update its priority or present value
func (d *GormDatabase) PointWriteByName(networkName, deviceName, pointName string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	getPnt, err := d.GetPointByName(networkName, deviceName, pointName)
	if err != nil {
		return nil, err
	}
	write, err := d.PointWrite(getPnt.UUID, body, fromPlugin)
	if err != nil {
		return nil, err
	}
	return write, nil
}

// CreatePoint creates an object.
func (d *GormDatabase) CreatePoint(body *model.Point, streamUUID string) (*model.Point, error) {
	var deviceModel *model.Device
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Point)
	deviceUUID := body.DeviceUUID
	body.Name = nameIsNil(body.Name)
	existing := d.pointNameExists(body, body)
	if existing {
		eMsg := fmt.Sprintf("a point with existing name: %s exists", body.Name)
		return nil, errors.New(eMsg)
	}

	obj, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	body.ObjectType = string(obj)
	query := d.DB.Where("uuid = ? ", deviceUUID).First(&deviceModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Description == "" {
		body.Description = "na"
	}
	body.ThingClass = model.ThingClass.Point
	body.CommonEnable.Enable = utils.NewTrue()
	body.CommonFault.InFault = true
	body.CommonFault.MessageLevel = model.MessageLevel.NoneCritical
	body.CommonFault.MessageCode = model.CommonFaultCode.PluginNotEnabled
	body.CommonFault.Message = model.CommonFaultMessage.PluginNotEnabled
	body.CommonFault.LastFail = time.Now().UTC()
	body.CommonFault.LastOk = time.Now().UTC()
	body.InSync = utils.NewFalse()
	body.WriteValueOnceSync = utils.NewFalse()
	if body.Priority == nil {
		body.Priority = &model.Priority{}
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, query.Error
	}
	if streamUUID != "" {
		producerModel := new(model.Producer)
		producerModel.StreamUUID = streamUUID
		producerModel.ProducerThingUUID = body.UUID
		producerModel.ProducerThingClass = model.ThingClass.Point
		producerModel.ProducerThingType = model.ThingClass.Point
		_, err := d.CreateProducer(producerModel)
		if err != nil {
			return nil, errors.New("ERROR on create new producer to an existing stream")
		}
	}
	plug, err := d.GetPluginIDFromDevice(deviceUUID)
	if err != nil {
		return nil, errors.New("ERROR failed to get plugin uuid")
	}
	t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, plug.PluginConfId, body.UUID)
	d.Bus.RegisterTopic(t)
	err = d.Bus.Emit(eventbus.CTX(), t, body)
	if err != nil {
		return nil, errors.New("ERROR on device eventbus")
	}
	return body, query.Error
}

// UpdatePoint update it.
func (d *GormDatabase) UpdatePoint(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Name != "" {
		existing := d.pointNameExists(pointModel, body)
		if existing {
			eMsg := fmt.Sprintf("a point with existing name: %s exists", body.Name)
			return nil, errors.New(eMsg)
		}
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&pointModel, body.Tags); err != nil {
			return nil, err
		}
	}
	body.InSync = utils.NewFalse()
	body.WriteValueOnceSync = utils.NewFalse()
	_ = d.DB.Model(&pointModel).Updates(&body)
	pnt, err := d.UpdatePointValue(uuid, pointModel, false)
	if err != nil {
		return nil, err
	}
	if !fromPlugin { //stop looping
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
	return pnt, nil
}

// PointWrite returns the device for the given id or nil.
func (d *GormDatabase) PointWrite(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Priority != nil {
		priority := map[string]interface{}{}
		priorityValue := reflect.ValueOf(*body.Priority)
		typeOfPriority := priorityValue.Type()
		highestPri := utils.NewArray()
		highestValue := utils.NewMap()
		for i := 0; i < priorityValue.NumField(); i++ {
			if priorityValue.Field(i).Type().Kind().String() == "ptr" {
				val := priorityValue.Field(i).Interface().(*float64)
				if val == nil {
					priority[typeOfPriority.Field(i).Name] = nil
				} else {
					highestPri.Add(i)
					highestValue.Set(i, *val)
					priority[typeOfPriority.Field(i).Name] = *val
				}
			}
		}
		notNil := false
		for _, v := range priority { //check if there is a value in priority array
			if v != nil {
				notNil = true
			}
		}
		if notNil {
			min, _ := highestPri.MinMaxInt() //get the highest priority
			val := highestValue.Get(min)     //get the highest priority value
			body.CurrentPriority = &min      //TODO check conversion
			v := val.(float64)
			body.PresentValue = &v //process the units as in temperature conversion
		}
		d.DB.Model(&pointModel.Priority).Updates(&priority)
	}
	var pntSync model.Point
	pntSync.WriteValueOnceSync = utils.NewFalse()
	point, err := d.UpdatePoint(uuid, &pntSync, fromPlugin)
	if err != nil {
		return nil, err
	}
	if !fromPlugin { //stop looping
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
	return point, nil
}

// UpdatePointValue returns the device for the given id or nil.
func (d *GormDatabase) UpdatePointValue(uuid string, body *model.Point, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	//TODO add in a check to make sure user doesn't set the addressID and the ObjectType the same as another point
	//check if there is an existing device with this address code
	_, err := checkObjectType(body.ObjectType)
	if err != nil {
		return nil, err
	}
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").Find(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	var presentValue *float64

	var limitMin *float64
	if body.LimitMin == nil {
		if pointModel.LimitMin != nil {
			limitMin = pointModel.LimitMin
		} else {
			limitMin = body.LimitMin
		}
	} else {
		limitMin = body.LimitMin
	}
	var limitMax *float64
	if body.LimitMax == nil {
		limitMax = pointModel.LimitMax
	}
	var scaleInMin *float64
	if body.ScaleInMin == nil {
		scaleInMin = pointModel.ScaleInMin
	}
	var scaleInMax *float64
	if body.ScaleInMax == nil {
		scaleInMax = pointModel.ScaleInMax
	}
	var scaleOutMin *float64
	if body.ScaleOutMin == nil {
		scaleOutMin = pointModel.ScaleOutMin
	}
	var scaleOutMax *float64
	if body.ScaleOutMax == nil {
		scaleOutMax = pointModel.ScaleOutMax
	}
	if body.Priority != nil {
		priority := map[string]interface{}{}
		priorityValue := reflect.ValueOf(*body.Priority)
		typeOfPriority := priorityValue.Type()
		highestPri := utils.NewArray()
		highestValue := utils.NewMap()
		for i := 0; i < priorityValue.NumField(); i++ {
			if priorityValue.Field(i).Type().Kind().String() == "ptr" {
				val := priorityValue.Field(i).Interface().(*float64)
				if val == nil {
					priority[typeOfPriority.Field(i).Name] = nil
				} else {
					log.Debug("UpdatePointValue() - i, val: ", typeOfPriority.Field(i).Name, *val)
					highestPri.Add(i)
					highestValue.Set(i, *val)
					priority[typeOfPriority.Field(i).Name] = *val
				}
			}
		}
		notNil := false
		for _, v := range priority { //check if there is a value in priority array
			if v != nil {
				notNil = true
			}
		}
		log.Debug("UpdatePointValue() - body.Fallback: ", body.Fallback)
		if notNil {
			min, _ := highestPri.MinMaxInt() //get the highest priority
			val := highestValue.Get(min)     //get the highest priority value
			body.CurrentPriority = &min      //TODO check conversion
			v := val.(float64)
			presentValue = &v //process the units as in temperature conversion
		} else if !utils.FloatIsNilCheck(body.Fallback) {
			body.Priority.P16 = utils.NewFloat64(*body.Fallback)
			body.OriginalValue = utils.NewFloat64(*body.Fallback)
			body.CurrentPriority = utils.NewInt(16)
			presentValue = utils.NewFloat64(*body.Fallback)
			log.Debug("UpdatePointValue() - FloatIsNilCheck(body.Fallback):", body.Priority.P16, body.OriginalValue, body.CurrentPriority, presentValue)
		}
		d.DB.Model(&pointModel.Priority).Updates(&priority)
	}
	ov := utils.Float64IsNil(presentValue)
	body.OriginalValue = &ov
	if presentValue != nil {
		log.Debug("UpdatePointValue() - *presentValue1: ", *presentValue)
	}
	if scaleInMin != nil {
		log.Debug("UpdatePointValue() - *scaleInMin: ", *scaleInMin)
	}
	if scaleInMax != nil {
		log.Debug("UpdatePointValue() - *scaleInMax: ", *scaleInMax)
	}
	if scaleOutMin != nil {
		log.Debug("UpdatePointValue() - *scaleOutMin: ", *scaleOutMin)
	}
	if scaleOutMax != nil {
		log.Debug("UpdatePointValue() - *scaleOutMax:", *scaleOutMax)
	}
	presentValue = pointScale(presentValue, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax)
	log.Debug("UpdatePointValue() - presentValue2: ", *presentValue)
	presentValue = pointRange(presentValue, limitMin, limitMax)
	log.Debug("UpdatePointValue() - presentValue3: ", *presentValue)
	eval, err := pointEval(presentValue, body.OriginalValue, pointModel.EvalMode, pointModel.Eval)
	if err != nil {
		log.Errorf("ERROR on point invalid point unit")
		return nil, err
	} else {
		pointModel.PresentValue = eval
	}
	vv, ok, err := pointUnits(pointModel)
	if err != nil {
		log.Errorf("ERROR on point invalid point unit")
		return nil, err
	}
	if ok {
		presentValue = &vv
	}
	if !utils.Unit32NilCheck(pointModel.Decimal) {
		val := utils.RoundTo(*presentValue, *pointModel.Decimal)
		presentValue = &val
	}
	if pointModel.IsProducer == nil && body.IsProducer == nil {
		if compare(pointModel, body) {
			_, err := d.ProducerWrite("point", pointModel)
			if err != nil {
				log.Errorf("ERROR ProducerPointCOV at func UpdatePointByFieldAndType")
				return nil, err
			}
		}
	}
	body.PresentValue = presentValue
	_ = d.DB.Model(&pointModel).Updates(&body)
	if !fromPlugin { //stop looping
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

// GetPointByField returns the point for the given field ie name or nil.
func (d *GormDatabase) GetPointByField(field string, value string) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? ", field)
	query := d.DB.Where(f, value).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

// PointAndQuery will do an SQL AND
func (d *GormDatabase) PointAndQuery(value1 string, value2 string) (*model.Point, error) {
	var pointModel *model.Point
	f := "object_type = ? AND address_id = ?"
	query := d.DB.Where(f, value1, value2).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

// GetPointByFieldAndIOID get by field and update.
func (d *GormDatabase) GetPointByFieldAndIOID(field string, value string, body *model.Point) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? AND io_id = ?", field)
	query := d.DB.Where(f, value, body.IoID).First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

// GetPointByFieldAndThingType get by field and thing_type.
func (d *GormDatabase) GetPointByFieldAndThingType(field string, value string, body *model.Point) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? AND thing_type = ?", field)
	query := d.DB.Where(f, value, body.ThingType).First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

// UpdatePointByFieldAndUnit get by field and update.
func (d *GormDatabase) UpdatePointByFieldAndUnit(field string, value string, body *model.Point, writeValue bool) (*model.Point, error) {
	var pointModel *model.Point
	f := fmt.Sprintf("%s = ? AND unit_type = ?", field)
	query := d.DB.Where(f, value, body.UnitType).First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	_, err := d.UpdatePointValue(pointModel.UUID, body, true)
	if err != nil {
		return nil, err
	}
	if pointModel.IsProducer == nil {
		if compare(pointModel, body) {
			_, err := d.ProducerWrite("point", pointModel)
			if err != nil {
				log.Errorf("ERROR ProducerPointCOV at func UpdatePointByFieldAndType")
				return nil, err
			}
		}
	}
	return pointModel, nil
}

// DeletePoint delete a Device.
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

// DropPoints delete all points.
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
