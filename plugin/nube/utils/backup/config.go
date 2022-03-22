package main

type Config struct {
	Host            string `yaml:"host"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	BucketName      string `yaml:"bucket_name"`
}

func (i *Instance) DefaultConfig() interface{} {
	return &Config{
		Host:            "127.0.0.1:9000",
		AccessKeyID:     "",
		SecretAccessKey: "",
		BucketName:      "",
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
