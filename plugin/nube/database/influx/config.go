package main

type Influx struct {
	Host        string  `yaml:"host"`
	Port        int     `yaml:"port"`
	Token       *string `yaml:"token"`
	Org         string  `yaml:"org"`
	Bucket      string  `yaml:"bucket"`
	Measurement string  `yaml:"measurement"`
}

type Job struct {
	Frequency string `yaml:"frequency"`
}

type Config struct {
	Influx []Influx `yaml:"influx"`
	Job    Job      `yaml:"job"`
}

func (i *Instance) DefaultConfig() interface{} {
	influx := Influx{
		Host:        "localhost",
		Port:        8086,
		Token:       nil,
		Org:         "nube-org",
		Bucket:      "nube-bucket",
		Measurement: "points",
	}
	job := Job{
		Frequency: "1m",
	}
	return &Config{
		Influx: []Influx{influx},
		Job:    job,
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
