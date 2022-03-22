package main

func (inst *Instance) setUUID() {
	q, err := inst.db.GetPluginByPath(name)
	if err != nil {
		return
	}
	inst.pluginUUID = q.UUID
}
