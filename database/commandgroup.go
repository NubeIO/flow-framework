package database

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/flow-framework/utils"
)

func (d *GormDatabase) GetCommandGroups() ([]*model.CommandGroup, error) {
	var commandGroup []*model.CommandGroup
	query := d.DB.Find(&commandGroup)
	if query.Error != nil {
		return nil, query.Error
	}
	return commandGroup, nil
}

func (d *GormDatabase) GetCommandGroup(uuid string) (*model.CommandGroup, error) {
	var commandGroup *model.CommandGroup
	query := d.DB.Where("uuid = ? ", uuid).First(&commandGroup)
	if query.Error != nil {
		return nil, query.Error
	}
	return commandGroup, nil
}

func (d *GormDatabase) CreateCommandGroup(body *model.CommandGroup) (*model.CommandGroup, error) {
	body.UUID = utils.MakeTopicUUID(model.CommonNaming.CommandGroup)
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

func (d *GormDatabase) UpdateCommandGroup(uuid string, body *model.CommandGroup) (*model.CommandGroup, error) {
	var commandGroup *model.CommandGroup
	query := d.DB.Where("uuid = ?", uuid).First(&commandGroup)
	if query.Error != nil {
		return nil, query.Error
	}
	query = d.DB.Model(&commandGroup).Updates(body)
	if query.Error != nil {
		return nil, query.Error
	}
	return commandGroup, nil
}

func (d *GormDatabase) DeleteCommandGroup(uuid string) (bool, error) {
	var commandGroup *model.CommandGroup
	query := d.DB.Where("uuid = ? ", uuid).Delete(&commandGroup)
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

func (d *GormDatabase) DropCommandGroups() (bool, error) {
	var commandGroup *model.CommandGroup
	query := d.DB.Where("1 = 1").Delete(&commandGroup)
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
