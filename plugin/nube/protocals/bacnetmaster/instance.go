package main

func (i *Instance) setUUID() {
	q, err := i.db.GetPluginByPath(path)
	if err != nil {
		return
	}
	i.pluginUUID = q.UUID
}
