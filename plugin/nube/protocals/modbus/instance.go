package main

func (inst *Instance) setUUID() {
	q, err := inst.db.GetPluginByPath(path)
	if err != nil {
		return
	}
	inst.pluginUUID = q.UUID
}
