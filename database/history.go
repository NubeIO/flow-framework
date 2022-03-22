package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/model"
	log "github.com/sirupsen/logrus"
)

// GetHistories returns all histories.
func (d *GormDatabase) GetHistories(args api.Args) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.buildHistoryQuery(args)
	query.Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// GetHistoriesForSync returns all histories after id.
// We order by `uuid` i.e. `producer_uuid`, so all similar data comes on same block which helps to reduce data query
// for fetching the data from points, devices, networks etc.
func (d *GormDatabase) GetHistoriesForSync(lastSyncId int) ([]*model.History, error) {
	var historiesModel []*model.History
	query := d.DB.Where("id > (?)", lastSyncId)
	query.Order("uuid ASC").Find(&historiesModel)
	if query.Error != nil {
		return nil, query.Error
	}
	return historiesModel, nil
}

// GetHistoriesByUUID returns the history for the given uuid or nil.
func (d *GormDatabase) GetHistoriesByUUID(uuid string, args api.Args) ([]*model.History, int64, error) {
	var count int64
	var historiesModel []*model.History
	query := d.buildHistoryQuery(args)
	query = query.Where("uuid = ?", uuid)
	query.Find(&historiesModel)
	query.Count(&count)
	return historiesModel, count, nil
}

// CreateHistory creates a thing.
func (d *GormDatabase) CreateHistory(body *model.History) (*model.History, error) {
	if err := d.DB.Create(&body).Error; err != nil {
		return nil, err
	}
	return body, nil
}

// CreateBulkHistory bulk creates a thing.
func (d *GormDatabase) CreateBulkHistory(histories []*model.History) (bool, error) {
	tx := d.DB.Begin() // for restricting the access by data source while bulk history creation is still to complete
	for _, history := range histories {
		_, err := d.CreateHistory(history)
		if err != nil {
			log.Error(fmt.Sprintf("Issue on creating history.id = %d, producer_uuid = %s", history.ID, history.UUID))
		}
	}
	tx.Commit()
	return true, nil
}

// DeleteHistoriesByUUID delete all history for given uuid.
func (d *GormDatabase) DeleteHistoriesByUUID(uuid string, args api.Args) (bool, error) {
	var historyModel *model.History
	query := d.buildHistoryQuery(args)
	query = query.Where("uuid = ? ", uuid)
	query.Delete(&historyModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	}
	return true, nil
}

// DeleteHistory delete a history for given id
func (d *GormDatabase) DeleteHistory(id int) (bool, error) {
	var historyModel *model.History
	query := d.DB.Where("id = ? ", id).Delete(&historyModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	}
	return true, nil
}

// DeleteAllHistories delete all histories.
func (d *GormDatabase) DeleteAllHistories() (bool, error) {
	var historyModel *model.History
	query := d.DB.Where("1 = 1")
	query.Delete(&historyModel)
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, nil
	}
	return true, nil
}
