package dbhandler

import "github.com/NubeIO/flow-framework/model"

func (h *Handler) GetEnabledIntegrationByPluginConfId(pcId string) ([]*model.Integration, error) {
	return getDb().GetEnabledIntegrationByPluginConfId(pcId)
}
