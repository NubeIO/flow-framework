package database

import (
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

type Job struct {
	*model.Job
}

func (d *GormDatabase) GetJobs() ([]*model.Job, error) {
	var jobsModel []*model.Job
	query := d.DB.Find(&jobsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobsModel, nil
}

func (d *GormDatabase) CreateJob(body *model.Job) (*model.Job, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.Job)
	body.Name = nameIsNil(body.Name)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) GetJob(uuid string) (*model.Job, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ? ", uuid).First(&jobModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil
}

func (d *GormDatabase) GetJobsByPluginConfigId(pcId string) ([]*model.Job, error) {
	var jobsModel []*model.Job
	query := d.DB.Where("plugin_conf_id = ?", pcId).Find(&jobsModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobsModel, nil
}

func (d *GormDatabase) GetJobByPluginConfId(pcId string) (*model.Job, error) {
	var jobModel *model.Job
	query := d.DB.Where("plugin_conf_id = ?", pcId).First(&jobModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil
}

// DeleteJob delete a job
func (d *GormDatabase) DeleteJob(uuid string) (bool, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ? ", uuid).Delete(&jobModel)
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

// UpdateJob  returns the device for the given id or nil.
func (d *GormDatabase) UpdateJob(uuid string, body *model.Job) (*model.Job, error) {
	var jobModel *model.Job
	query := d.DB.Where("uuid = ?", uuid).First(&jobModel)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&jobModel).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return jobModel, nil

}

// DropJobs delete all.
func (d *GormDatabase) DropJobs() (bool, error) {
	var jobModel *model.Job
	query := d.DB.Where("1 = 1").Delete(&jobModel)
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
