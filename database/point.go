package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/flow-framework/utils/float"
	"github.com/NubeIO/flow-framework/utils/integer"
	"github.com/NubeIO/flow-framework/utils/nmath"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/flow-framework/utils/priorityarray"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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
	body.UUID = nuuid.MakeTopicUUID(model.ThingClass.Point)
	body.Name = nameIsNil(body.Name)

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
		body.PointPriorityArrayMode = model.PriorityArrayToPresentValue // sets default priority array mode.
	}
	body.ThingClass = model.ThingClass.Point
	body.CommonEnable.Enable = boolean.NewTrue()
	body.InSync = boolean.NewFalse()
	if body.Priority == nil {
		body.Priority = &model.Priority{}
	}
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	network, err := d.GetNetworkByDeviceUUID(body.DeviceUUID, api.Args{})
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
	// check for mapping
	if network.AutoMappingNetworksSelection != "" {
		pointMapping := &model.PointMapping{}
		pointMapping.Point = body
		pointMapping.AutoMappingFlowNetworkName = network.AutoMappingFlowNetworkName
		pointMapping.AutoMappingFlowNetworkUUID = network.AutoMappingFlowNetworkUUID
		pointMapping.AutoMappingNetworksSelection = []string{network.AutoMappingNetworksSelection}
		pointMapping.AutoMappingEnableHistories = network.AutoMappingEnableHistories
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

	if !integer.IsNil(body.AddressID) {
		if existingAddrID {
			eMsg := fmt.Sprintf("a point with existing AddressID: %d exists", integer.NonNil(body.AddressID))
			return nil, errors.New(eMsg)
		}
	}
	if len(body.Tags) > 0 {
		if err := d.updateTags(&pointModel, body.Tags); err != nil {
			return nil, err
		}
	}
	// example modbus: if user changes the data type then do a new read of the point on the modbus network
	if !fromPlugin {
		pointModel.InSync = boolean.NewFalse()
	}
	// TODO: ARE THESE REQUIRED? OR ARE THEY DONE WITH THE FOLLOWING DB CALL?
	pointModel.WritePollRequired = body.WritePollRequired
	pointModel.ReadPollRequired = body.ReadPollRequired
	pointModel.WriteMode = body.WriteMode
	pointModel.PollPriority = body.PollPriority
	pointModel.PollRate = body.PollRate

	if body.Fallback == nil {
		// nil is ignored on GORM, so we are pushing forcefully because fallback is nullable field
		d.DB.Model(&pointModel).Update("fallback", nil)
	}
	query = d.DB.Model(&pointModel).Updates(&body)

	//TODO: we need to decide if a read only point needs to have a priority array or if it should just be nil.
	if body.Priority == nil {
		pointModel.Priority = &model.Priority{}
	} else {
		pointModel.Priority = body.Priority
	}

	priorityMap := priorityarray.ConvertToMap(*pointModel.Priority)
	pnt, _, _, _, err := d.UpdatePointValue(pointModel, &priorityMap, fromPlugin)
	if err != nil {
		return nil, err
	}
	return pnt, nil
}

func (d *GormDatabase) PointWrite(uuid string, body *model.PointWriter, fromPlugin bool) (returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	var pointModel *model.Point
	query := d.DB.Where("uuid = ?", uuid).Preload("Priority").First(&pointModel)
	if query.Error != nil {
		return nil, false, false, false, query.Error
	}
	if body.Priority == nil {
		return nil, false, false, false, errors.New("no priority value is been sent")
	} else {
		pointModel.ValueUpdatedFlag = boolean.NewTrue()
	}
	pointModel.InSync = boolean.NewFalse()
	pointModel.WritePollRequired = boolean.NewTrue()
	point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err := d.UpdatePointValue(pointModel, body.Priority, fromPlugin)
	return point, isPresentValueChange, isWriteValueChange, isPriorityChanged, err
}

func (d *GormDatabase) UpdatePointValue(pointModel *model.Point, priority *map[string]*float64, fromPlugin bool) (returnPoint *model.Point, isPresentValueChange, isWriteValueChange, isPriorityChanged bool, err error) {
	if pointModel.PointPriorityArrayMode == "" {
		pointModel.PointPriorityArrayMode = model.PriorityArrayToPresentValue // sets default priority array mode
	}

	pointModel, priority, presentValue, writeValue, isPriorityChanged := d.updatePriority(pointModel, priority)
	ov := float.Copy(presentValue)
	pointModel.OriginalValue = ov

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
	// example for wires and modbus:
	// if a new value is written from wires then set this to false so the modbus knows on the next poll to write a new
	// value to the modbus point
	if !fromPlugin {
		pointModel.InSync = boolean.NewFalse()
	}
	if !integer.IsUnit32Nil(pointModel.Decimal) && presentValue != nil {
		val := nmath.RoundTo(*presentValue, *pointModel.Decimal)
		presentValue = &val
	}

	isPresentValueChange = !float.ComparePtrValues(pointModel.PresentValue, presentValue) // Use for createCOVHistory
	// isWriteValueChange
	// in some cases we don't write presentValue directly to the model.
	// Instead of writing that, we firstly write it on writeValue & and we read that writeValue to update presentValue.
	// Examples are: modbus, edge28 plugins
	// So for such cases, to trigger that value we do this comparison
	isWriteValueChange = !float.ComparePtrValues(pointModel.WriteValue, writeValue)
	isChange := isPresentValueChange || isWriteValueChange || isPriorityChanged || presentValueTransformFault
	pointModel.WriteValue = writeValue
	if writeValue == nil {
		// nil is ignored on GORM, so we are pushing forcefully because isChange comparison will fail on `null` write
		// if we don't do this: it goes on infinite looping as writeValue and point.WriteValue gets different values
		d.DB.Model(&pointModel).Update("write_value", nil)
	}

	// If the present value transformations have resulted in an error, DB needs to be updated with the errors,
	// but PresentValue should not change
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

	if isChange {
		_ = d.DB.Model(&pointModel).Updates(&pointModel)
		err = d.ProducersPointWrite(pointModel.UUID, priority, pointModel.PresentValue, isPresentValueChange)
		if err != nil {
			return nil, false, false, false, err
		}
		d.DB.Model(&model.Writer{}).
			Where("writer_thing_uuid = ?", pointModel.UUID).
			Update("present_value", pointModel.PresentValue)
	}
	if !fromPlugin { // stop looping
		plug, err := d.GetNetworkByDeviceUUID(pointModel.DeviceUUID, api.Args{})
		if err != nil {
			return nil, false, false, false, errors.New("ERROR failed to get plugin uuid")
		}
		t := fmt.Sprintf("%s.%s.%s", eventbus.PluginsUpdated, plug.PluginConfId, pointModel.UUID)
		d.Bus.RegisterTopic(t)
		err = d.Bus.Emit(eventbus.CTX(), t, pointModel)
		if err != nil {
			return nil, false, false, false, errors.New("ERROR on device eventbus")
		}
	}
	return pointModel, isPresentValueChange, isWriteValueChange, isPriorityChanged, nil
}

func (d *GormDatabase) DeletePoint(uuid string) (bool, error) {
	point, err := d.GetPoint(uuid, api.Args{})
	if err != nil {
		return false, err
	}
	producers, _ := d.GetProducers(api.Args{ProducerThingUUID: &point.UUID})
	if producers != nil {
		for _, producer := range producers {
			_, _ = d.DeleteProducer(producer.UUID)
		}
	}
	writers, _ := d.GetWriters(api.Args{WriterThingUUID: &point.UUID})
	if writers != nil {
		for _, writer := range writers {
			_, _ = d.DeleteWriter(writer.UUID)
		}
	}
	query := d.DB.Delete(&point)
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
