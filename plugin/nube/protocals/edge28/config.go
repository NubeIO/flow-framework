package main

import "github.com/NubeIO/flow-framework/plugin/nube/protocals/edge28/config"

func (inst *Instance) DefaultConfig() interface{} {
	return &config.Config{
		EnablePolling: true,
		LogLevel:      "DEBUG", // DEBUG or ERROR
		PollRate:      5,       // in seconds
	}
}

func (inst *Instance) GetConfig() interface{} {
	return inst.config
}

func (inst *Instance) ValidateAndSetConfig(c interface{}) error {
	newConfig := c.(*config.Config)
	inst.config = newConfig
	return nil
}
