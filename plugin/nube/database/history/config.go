package main

type Job struct {
	Frequency string `yaml:"frequency"`
}

type Config struct {
	Job Job `yaml:"job"`
}

func (i *Instance) DefaultConfig() interface{} {
	job := Job{
		Frequency: "1m",
	}
	return &Config{
		Job: job,
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
