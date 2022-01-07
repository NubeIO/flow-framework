package dbhandler

import "github.com/NubeIO/flow-framework/model"

func (h *Handler) GetLatestProducerHistoryByProducerName(name string) (*model.ProducerHistory, error) {
	q, err := getDb().GetLatestProducerHistoryByProducerName(name)
	if err != nil {
		return nil, err
	}
	return q, nil
}
