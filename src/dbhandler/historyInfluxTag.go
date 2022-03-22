package dbhandler

import "github.com/NubeIO/flow-framework/model"

func (h *Handler) GetHistoryInfluxTags(producerUuid string) ([]*model.HistoryInfluxTag, error) {
	return getDb().GetHistoryInfluxTags(producerUuid)
}
