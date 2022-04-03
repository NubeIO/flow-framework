package plugin

type dbStorageHandler struct {
	pluginID string
	db       Database
}

func (c dbStorageHandler) Save(b []byte) error {
	conf, err := c.db.GetPlugin(c.pluginID)
	if err != nil {
		return err
	}
	conf.Storage = b
	return c.db.UpdatePluginConf(conf)
}

func (c dbStorageHandler) Load() ([]byte, error) {
	pluginConf, err := c.db.GetPlugin(c.pluginID)
	if err != nil {
		return nil, err
	}
	return pluginConf.Storage, nil
}
