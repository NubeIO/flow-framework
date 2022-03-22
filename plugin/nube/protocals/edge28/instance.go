package main

func (i *Instance) setUUID() {
	q, err := i.db.GetPluginByPath(name)
	if err != nil {
		return
	}
	i.pluginUUID = q.UUID
}
