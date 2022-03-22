package main

type Config struct {
	MagicString string `yaml:"magic_string"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		MagicString: "N/A",
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	inst.config = newConfig
	return nil
}
