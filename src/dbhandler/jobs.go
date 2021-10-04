package dbhandler

import (
	"github.com/NubeDev/flow-framework/model"
)

func (h *Handler) GetJobs() ([]*model.Job, error) {
	q, err := getDb().GetJobs()
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetJob(uuid string) (*model.Job, error) {
	q, err := getDb().GetJob(uuid)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) GetJobByPluginConfId(pcId string) (*model.Job, error) {
	q, err := getDb().GetJobByPluginConfId(pcId)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) CreateJob(body *model.Job) (*model.Job, error) {
	q, err := getDb().CreateJob(body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) UpdateJob(uuid string, body *model.Job) (*model.Job, error) {
	q, err := getDb().UpdateJob(uuid, body)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (h *Handler) DeleteJob(uuid string) (bool, error) {
	_, err := getDb().DeleteJob(uuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
