package main

func (inst *Instance) setUUID() {
	q, err := inst.db.GetPluginByPath(path)
	if err != nil {
		return
	}
	inst.pluginUUID = q.UUID

	//nrest_bacnet_server.NewClient(rt)
	//aa := nrest_bacnet_server.RestClient{}
	//nrest_bacnet_server.BacnetPoint{}
}
