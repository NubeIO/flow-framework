package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func (d *GormDatabase) GetPoints(args api.Args) ([]*model.Point, error) {
	var pointsModel []*model.Point
	query := d.buildPointQuery(args)
	if err := query.Find(&pointsModel).Error; err != nil {
		return nil, err
	}
	return pointsModel, nil
}

func (d *GormDatabase) GetPointsBulk(bulkPoints []*model.Point) ([]*model.Point, error) {
	var pointsModel []*model.Point
	points, err := d.GetPoints(api.Args{WithPriority: true})
	if err != nil {
		return nil, err
	}
	for _, pnt := range points {
		for _, search := range bulkPoints {
			if pnt.UUID == search.UUID {
				pointsModel = append(pointsModel, pnt)
			}

		}
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
	body.Name = nameIsNil(body.Name)
	//existingAddrID := false
	//existingName, _ := d.pointNameExists(body)
	if body.AddressID != nil {
		//_, existingAddrID = d.pointNameExists(body)
	}
	//if existingName {
	//	eMsg := fmt.Sprintf("a point with existing name: %s exists", body.Name)
	//	return nil, errors.New(eMsg)
	//}
	//if existingAddrID {
	//	eMsg := fmt.Sprintf("a point with existing AddressID: %d exists", utils.IntIsNil(body.AddressID))
	//	return nil, errors.New(eMsg)
	//}
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
	if body.PointPriorityArrayMode == "" {
		body.PointPriorityArrayMode = model.PriorityArrayToPresentValue //sets default priority array mode.
	}
	body.ThingClass = model.ThingClass.Point
	body.CommonEnable.Enable = utils.NewTrue()
	body.InSync = utils.NewFalse()
	if body.Priority == nil {
		body.Priority = &model.Priority{}
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	deviceUUID := body.DeviceUUID // TODO: Should there be a check to ensure that the DeviceUUID is sent with body?
	network, err := d.GetNetworkByDeviceUUID(deviceUUID, api.Args{})
	fmt.Printf("network: %+v\n", network)
	if err != nil {
		return nil, errors.New("ERROR failed to get plugin uuid")
	}
	if network == nil {
		return nil, errors.New("ERROR failed to get network")
	}
	if !fromPlugin {
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsCreated, network.PluginConfId, body.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, body)
		if err != nil {
			return nil, errors.New("ERROR on device eventbus")
		}
	}
	//check for mapping
	if network.AutoMappingNetworksSelection != "" {
		pointMapping := &model.PointMapping{}
		pointMapping.Point = body
		pointMapping.AutoMappingFlowNetworkName = network.AutoMappingFlowNetworkName
		pointMapping.AutoMappingFlowNetworkUUID = network.AutoMappingFlowNetworkUUID
		pointMapping.AutoMappingNetworksSelection = []string{network.AutoMappingNetworksSelection}
		pointMapping, err = d.CreatePointMapping(pointMapping)
		if err != nil {
			log.Errorln("points.db.CreatePoint() failed to make auto point mapping")
			return nil, err
		} else {
			log.Println("points.db.CreatePoint() added point new mapping")
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
	// TODO: ARE THESE REQUIRED? OR ARE THEY DONE WITH THE FOLLOWING DB CALL?
	pointModel.WritePollRequired = body.WritePollRequired
	pointModel.ReadPollRequired = body.ReadPollRequired
	pointModel.WriteMode = body.WriteMode
	pointModel.PollPriority = body.PollPriority
	pointModel.PollRate = body.PollRate

	query = d.DB.Model(&pointModel).Updates(&body)
	// Don't update point value if priority array on body is nil
	//if body.Priority == nil {
	//	return pointModel, nil
	//} else {
	//	pointModel.Priority = body.Priority
	//}
	//
	//pnt, err := d.UpdatePointValue(pointModel, fromPlugin)
	//if err != nil {
	//	return nil, err
	//}
	//return pnt, nil
	return pointModel, nil
}

func (d *GormDatabase) PointWrite(uuid string, body *model.PointWriter, fromPlugin bool) (*model.Point, error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, query.Error
	}
	if body.Priority == nil {
		return nil, errors.New("no priority value is been sent")
	} else {
		pointModel.ValueUpdatedFlag = utils.NewTrue()
	}
	pointModel.InSync = utils.NewFalse()
	pointModel.WritePollRequired = utils.NewTrue()
	point, err := d.UpdatePointValue(pointModel, body.Priority, fromPlugin)
	return point, err
}

func (d *GormDatabase) UpdatePointValue(pointModel *model.Point, priority *map[string]*float64, fromPlugin bool) (*model.Point, error) {
	if pointModel.PointPriorityArrayMode == "" {
		pointModel.PointPriorityArrayMode = model.PriorityArrayToPresentValue //sets default priority array mode.
	}

	pointModel, presentValue := d.updatePriority(pointModel)
	//presentValue will be OriginalValue if PointPriorityArrayMode is PriorityArrayToWriteValue or ReadOnlyNoPriorityArrayRequired. OriginalValue will have been updated by plugin on successful read.
	ov := utils.Float64IsNil(presentValue)
	pointModel.OriginalValue = &ov

	presentValueTransformFault := false
	presentValue = pointScale(presentValue, pointModel.ScaleInMin, pointModel.ScaleInMax, pointModel.ScaleOutMin, pointModel.ScaleOutMax)
	presentValue = pointRange(presentValue, pointModel.LimitMin, pointModel.LimitMax)
	eval, err := pointEval(presentValue, pointModel.MathOnPresentValue)
	if err != nil {
		log.Errorln("point.db UpdatePointValue() error on run point MathOnPresentValue error:", err)
		pointModel.CommonFault.InFault = true
		pointModel.CommonFault.MessageLevel = model.MessageLevel.Warning
		pointModel.CommonFault.MessageCode = model.CommonFaultCode.PointError
		pointModel.CommonFault.Message = fmt.Sprint("point.db UpdatePointValue() error on run point MathOnPresentValue error:", err)
		pointModel.CommonFault.LastFail = time.Now().UTC()
		presentValueTransformFault = true
	} else {
		presentValue = eval
	}
	val, err := pointUnits(presentValue, pointModel.Unit, pointModel.UnitTo)
	if err != nil {
		log.Errorln("ERROR on point invalid point unit")
		pointModel.CommonFault.InFault = true
		pointModel.CommonFault.MessageLevel = model.MessageLevel.Warning
		pointModel.CommonFault.MessageCode = model.CommonFaultCode.PointError
		pointModel.CommonFault.Message = fmt.Sprint("point.db UpdatePointValue() invalid point units. error:", err)
		pointModel.CommonFault.LastFail = time.Now().UTC()
		presentValueTransformFault = true
	} else {
		presentValue = val
	}
	//example for wires and modbus: if a new value is written from  wires then set this to false so the modbus knows on the next poll to write a new value to the modbus point
	if !fromPlugin {
		pointModel.InSync = utils.NewFalse()
	}
	if !utils.Unit32NilCheck(pointModel.Decimal) && presentValue != nil {
		val := utils.RoundTo(*presentValue, *pointModel.Decimal)
		presentValue = &val
	}

	//TODO: I think this whole section can be deleted; Unless the overhead of d.ProducersPointWrite(pointModel) is too high.
	dbUpdateRequired := false
	switch pointModel.PointPriorityArrayMode {
	case model.PriorityArrayToPresentValue:
		//Priority array has already been pushed to DB so only need to update if PresentValue has changed or if there is a value transform error.
		dbUpdateRequired = !utils.CompareFloatPtr(pointModel.PresentValue, presentValue) || presentValueTransformFault

	case model.PriorityArrayToWriteValue:
		// TODO: Currently there isn't a good comparison to decide if a DB update is required, so just do it.  This could be added, but it would require a bunch of rework to the above functions.
		dbUpdateRequired = true

	case model.ReadOnlyNoPriorityArrayRequired:
		dbUpdateRequired = !utils.CompareFloatPtr(pointModel.PresentValue, presentValue) || presentValueTransformFault
	}

	// If the present value tranformations have resulted in an error, DB needs to be updated with the errors, but PresentValue should not change.
	if !presentValueTransformFault {
		pointModel.PresentValue = presentValue
	}
	if presentValue == nil {
		// nil is ignored on GORM, so we are pushing forcefully because isChange comparison will fail on `null` write
		d.DB.Model(&pointModel).Update("present_value", nil)
		d.DB.Model(&model.Writer{}).
			Where("writer_thing_uuid = ?", pointModel.UUID).
			Update("present_value", nil)
	}

	if dbUpdateRequired {
		_ = d.DB.Model(&pointModel).Updates(&pointModel)
		err = d.ProducersPointWrite(pointModel.UUID, priority, pointModel.PresentValue)
		if err != nil {
			return nil, err
		}
		d.DB.Model(&model.Writer{}).
			Where("writer_thing_uuid = ?", pointModel.UUID).
			Update("present_value", pointModel.PresentValue)
	}
	if !fromPlugin { // stop looping
		plug, err := d.GetNetworkByDeviceUUID(pointModel.DeviceUUID, api.Args{})
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
		plug, err := d.GetNetworkByDeviceUUID(point.DeviceUUID, api.Args{})
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
