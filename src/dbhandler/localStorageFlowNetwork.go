package dbhandler

func (h *Handler) RefreshLocalStorageFlowToken() (*bool, error) {
	return getDb().RefreshLocalStorageFlowToken()
}
