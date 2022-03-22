package main

type Config struct {
	MagicString string `yaml:"magic_string"`
}

func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		MagicString: "N/A",
	}
}

func (i *Instance) GetConfig() interface{} {
	return i.config
}

func (i *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	i.config = newConfig
	return nil
}
