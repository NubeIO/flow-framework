package database

import (
	"encoding/json"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

func (d *GormDatabase) GetSchedules() ([]*model.Schedule, error) {
	var scheduleModel []*model.Schedule
	query := d.DB.Find(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetSchedule(uuid string) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	query := d.DB.Where("uuid = ? ", uuid).First(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return scheduleModel, nil
}

func (d *GormDatabase) GetOneScheduleByArgs(args api.Args) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	query := d.buildScheduleQuery(args)
	if err := query.First(&scheduleModel).Error; err != nil {
		return nil, query.Error
	}
	return scheduleModel, nil
}

func (d *GormDatabase) CreateSchedule(body *model.Schedule) (*model.Schedule, error) {
	body.UUID = utils.MakeTopicUUID(model.ThingClass.Schedule)
	body.Name = nameIsNil(body.Name)
	body.ThingClass = model.ThingClass.Schedule
	body.ThingType = model.ThingClass.Schedule
	validSchedule, err := d.validateSchedule(body)
	if err != nil {
		return nil, err
	}
	body.Schedule = validSchedule
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) validateSchedule(schedule *model.Schedule) ([]byte, error) {
	if schedule.Schedule == nil {
		return nil, nil
	}
	scheduleDataModel := new(model.ScheduleData)
	err := json.Unmarshal(schedule.Schedule, &scheduleDataModel)
	if err != nil {
		return nil, err
	}
	validSchedule, err := json.Marshal(scheduleDataModel)
	if err != nil {
		return nil, err
	}
	return validSchedule, nil
}

func (d *GormDatabase) UpdateSchedule(uuid string, body *model.Schedule) (*model.Schedule, error) {
	var scheduleModel *model.Schedule
	validSchedule, err := d.validateSchedule(body)
	if err != nil {
		return nil, err
	}
	body.Schedule = validSchedule
	body.ThingClass = model.ThingClass.Schedule
	body.ThingType = model.ThingClass.Schedule
	query := d.DB.Where("uuid = ?", uuid).Find(&scheduleModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&scheduleModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return scheduleModel, nil
}

func (d *GormDatabase) ScheduleWrite(uuid string, body *model.ScheduleData) error {
	scheduleData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	schedule := map[string]interface{}{}
	schedule["schedule"] = scheduleData
	err = d.DB.Model(model.Schedule{}).Where("uuid = ?", uuid).Updates(schedule).Error
	if err != nil {
		return err
	}
	return d.ProducersScheduleWrite(uuid, body)
}

func (d *GormDatabase) DeleteSchedule(uuid string) (bool, error) {
	var schModel *model.Schedule
	query := d.DB.Where("uuid = ? ", uuid).Delete(&schModel)
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

func (d *GormDatabase) DropSchedules() (bool, error) {
	var schModel *model.Schedule
	query := d.DB.Where("1 = 1").Delete(&schModel)
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
