package dbhandler

import (
	"github.com/NubeIO/flow-framework/model"
)

func (h *Handler) GetHistoriesForSync(lastSyncId int) ([]*model.History, error) {
	return getDb().GetHistoriesForSync(lastSyncId)
}

func (h *Handler) CreateBulkHistory(histories []*model.History) (bool, error) {
	return getDb().CreateBulkHistory(histories)
}
