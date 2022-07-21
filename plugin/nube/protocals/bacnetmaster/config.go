package main

type Config struct {
	MagicString   string `yaml:"magic_string"`
	EnablePolling bool   `yaml:"enable_polling"`
	LogLevel      string `yaml:"log_level"`
}

func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		EnablePolling: true,
		LogLevel:      "ERROR", // DEBUG or ERROR
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(c interface{}) error {
	newConfig := c.(*Config)
	inst.config = newConfig
	return nil
}
