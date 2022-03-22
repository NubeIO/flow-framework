package main

// Config is user plugin configuration
type Config struct {
	MagicString string `yaml:"na"`
}

// DefaultConfig implements plugin.Configurer
func (inst *Instance) DefaultConfig() interface{} {
	return &Config{
		MagicString: "",
	}
}

// ValidateAndSetConfig implements plugin.Configurer
func (inst *Instance) ValidateAndSetConfig(config interface{}) error {
	newConfig := config.(*Config)
	inst.config = newConfig
	return nil
}
