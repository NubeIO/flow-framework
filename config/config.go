package config

import (
	"flag"
	"path"

	"github.com/NubeIO/configor"
	"github.com/NubeIO/flow-framework/utils/file"
	"github.com/NubeIO/flow-framework/utils/security"
)

type Configuration struct {
	Server struct {
		KeepAlivePeriodSeconds int
		ListenAddr             string `default:"0.0.0.0"`
		Port                   int
		RSPort                 int `default:"1616"`
		ResponseHeaders        map[string]string
		Stream                 struct {
			PingPeriodSeconds int `default:"45"`
			AllowedOrigins    []string
		}
		Cors struct {
			AllowOrigins []string
			AllowMethods []string
			AllowHeaders []string
		}
	}
	Database struct {
		Dialect    string `default:"sqlite3"`
		Connection string `default:"data.db"`
		LogLevel   string `default:"WARN"`
	}
	PassStrength int `default:"10"`
	LogLevel     string
	Location     struct {
		GlobalDir string `default:"./"`
		ConfigDir string `default:"config"`
		DataDir   string `default:"data"`
		Data      struct {
			PluginsDir        string `default:"plugins"`
			UploadedImagesDir string `default:"images"`
		}
		DeviceInfoFile    string `default:"/data/rubix-registry/device_info.json"`
		TokenFolder       string `default:"/data/rubix-service/data/"`
		InternalTokenFile string `default:"internal_token.txt"`
	}
	Prod            bool  `default:"false"`
	Auth            *bool `default:"true"`
	ProducerHistory struct {
		Cleaner struct {
			Enable              *bool `default:"true"`
			Frequency           int   `default:"600"`
			DataPersistingHours int   `default:"24"`
		}
		SyncInterval struct {
			Enable     *bool `default:"true"`
			SyncPeriod int   `default:"10"`
		}
	}
	SecretKey string
	MQTT      struct {
		Address string `default:"localhost"`
		Port    int    `default:"1883"`
	}
}

var config *Configuration = nil

func Get() *Configuration {
	return config
}

func CreateApp() *Configuration {
	config = new(Configuration)
	config = config.Parse()
	config = config.HandleSecretKey()
	err := configor.New(&configor.Config{EnvironmentPrefix: "FLOW"}).Load(config, path.Join(config.GetAbsConfigDir(), "config.yml"))
	if err != nil {
		panic(err)
	}
	return config
}

func (conf *Configuration) Parse() *Configuration {
	port := flag.Int("p", 1660, "Port")
	globalDir := flag.String("g", "./", "Global Directory")
	dataDir := flag.String("d", "data", "Data Directory")
	configDir := flag.String("c", "config", "Config Directory")
	prod := flag.Bool("prod", false, "Deployment Mode")
	mqttAddr := flag.String("mqtt-address", "localhost", "MQTT Broker Address")
	mqttPort := flag.Int("mqtt-port", 1883, "MQTT Broker Port")
	flag.Parse()
	conf.Server.Port = *port
	conf.Location.GlobalDir = *globalDir
	conf.Location.DataDir = *dataDir
	conf.Location.ConfigDir = *configDir
	conf.Prod = *prod
	conf.MQTT.Address = *mqttAddr
	conf.MQTT.Port = *mqttPort
	return conf
}

func (conf *Configuration) HandleSecretKey() *Configuration {
	secretKey, _ := file.ReadFile(path.Join(config.GetAbsConfigDir(), "secret.txt"))
	if secretKey == "" {
		secretKey = security.GenerateToken()
		_, _ = file.WriteDataToFileAsString(path.Join(config.GetAbsConfigDir(), "secret.txt"), secretKey)
	}
	conf.SecretKey = secretKey
	return conf
}

func (conf *Configuration) GetAbsDataDir() string {
	return path.Join(conf.Location.GlobalDir, conf.Location.DataDir)
}

func (conf *Configuration) GetAbsConfigDir() string {
	return path.Join(conf.Location.GlobalDir, conf.Location.ConfigDir)
}

func (conf *Configuration) GetAbsPluginDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.Location.Data.PluginsDir)
}

func (conf *Configuration) GetAbsUploadedImagesDir() string {
	return path.Join(conf.GetAbsDataDir(), conf.Location.Data.UploadedImagesDir)
}

func (conf *Configuration) GetAbsInternalTokenFile() string {
	return path.Join(conf.Location.TokenFolder, conf.Location.InternalTokenFile)
}
