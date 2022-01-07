package dbhandler

func (h *Handler) RefreshFlowNetworksConnections() (*bool, error) {
	return getDb().RefreshFlowNetworksConnections()
}
