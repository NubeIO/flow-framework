package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
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

func (d *GormDatabase) GetOnePointByArgs(args api.Args) (*model.Point, error) {
	var pointModel *model.Point
	query := d.buildPointQuery(args)
	if err := query.First(&pointModel).Error; err != nil {
		return nil, query.Error
	}
	return pointModel, nil
}

func (d *GormDatabase) CreatePoint(body *model.Point, fromPlugin bool) (*model.Point, error) {
	var deviceModel *model.Device
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Point)
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
	query := d.DB.Where("uuid = ? ", body.DeviceUUID).First(&deviceModel)
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
	var devArg api.Args
	dev, err := d.GetDevice(body.DeviceUUID, devArg)
	if err != nil {
		return nil, errors.New("ERROR failed to get device (for references)")
	}
	body.NetworkUUID = dev.NetworkUUID
	var netArg api.Args
	net, err := d.GetNetwork(dev.NetworkUUID, netArg)
	if err != nil {
		return nil, errors.New("ERROR failed to get network (for references)")
	}
	if !fromPlugin { //TODO: This looks like it should be reversed
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, net.PluginConfId, body.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, body)
		if err != nil {
			return nil, errors.New("ERROR on device eventbus")
		}
	}
	return body, query.Error
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
	eval, err := pointEval(presentValue, pointModel.OriginalValue, pointModel.EvalMode, pointModel.Eval)
	if err != nil {
		log.Errorf("ERROR on point invalid eval")
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
	}
	if isChange == true {
		_ = d.DB.Model(&pointModel).Updates(&pointModel)
		err = d.ProducersPointWrite(pointModel)
		if err != nil {
			return nil, err
		}
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

func (d *GormDatabase) updatePriority(pointModel *model.Point) (*model.Point, *float64) {
	var presentValue *float64
	if pointModel.Priority != nil {
		priorityMap, highestValue, currentPriority, isPriorityExist := d.parsePriority(pointModel.Priority)
		if isPriorityExist {
			pointModel.CurrentPriority = &currentPriority
			presentValue = &highestValue
		} else if !utils.FloatIsNilCheck(pointModel.Fallback) {
			pointModel.Priority.P16 = utils.NewFloat64(*pointModel.Fallback)
			pointModel.CurrentPriority = utils.NewInt(16)
			presentValue = utils.NewFloat64(*pointModel.Fallback)
		}
		d.DB.Model(&model.Priority{}).Where("point_uuid = ?", pointModel.UUID).Updates(&priorityMap)
	}
	return pointModel, presentValue
}

func (d *GormDatabase) parsePriority(priority *model.Priority) (map[string]interface{}, float64, int, bool) {
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
				}
				priorityMap[typeOfPriority.Field(i).Name] = *val
				isPriorityExist = true
			}
		}
	}
	return priorityMap, highestValue, currentPriority, isPriorityExist
}

func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	var pointModel *model.Point
	point, err := d.GetPoint(uuid, api.Args{})
	pntUUID := point.UUID
	devUUID := point.DeviceUUID
	netUUID := point.NetworkUUID
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
		if err != nil {
			return false, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s.%s", eventbus.PluginsDeleted, netUUID, pntUUID, devUUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, point)
		if err != nil {
			return false, errors.New("ERROR on point delete eventbus")
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
