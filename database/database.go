package database

import (
	"fmt"
	"github.com/NubeDev/flow-framework/auth/password"
	"github.com/NubeDev/flow-framework/cachestore"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/logger"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin"
	"github.com/NubeDev/flow-framework/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var mkdirAll = os.MkdirAll
var gormDatabase *GormDatabase

//var GetDatabaseBus eventbus.BusService

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, defaultUser, defaultPass string, strength int, logLevel string,
	createDefaultUserIfNotExist bool, production bool) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)
	_connection := fmt.Sprintf("%s?_foreign_keys=on", connection)
	db, err := gorm.Open(sqlite.Open(_connection), &gorm.Config{
		Logger: logger.New().SetLogMode(logLevel),
	})
	if err != nil {
		panic("failed to connect database")
	}
	var alerts []model.Alert
	var user []model.User
	var application []model.Application
	var message []model.Message
	var client []model.Client
	var pluginConf []model.PluginConf
	var network []model.Network
	var device []model.Device
	var point []model.Point
	var priority []model.Priority
	var producerHistory []model.ProducerHistory
	var consumerHistory []model.ConsumerHistory
	var flowNetwork []model.FlowNetwork
	var rubixPlat []model.RubixPlat
	var job []model.Job
	var stream []model.Stream
	var commandGroup []model.CommandGroup
	var producer []model.Producer
	var consumer []model.Consumer
	var writer []model.Writer
	var writerClone []model.WriterClone
	var integration []model.Integration
	var mqttConnection []model.MqttConnection
	var credentials []model.IntegrationCredentials
	var serialConnection []model.SerialConnection
	var ipConnection []model.IpConnection
	var schedule []model.Schedule
	var tags []model.Tag
	var models = []interface{}{
		&alerts,
		&user,
		&application,
		&message,
		&client,
		&pluginConf,
		&network,
		&device,
		&point,
		&rubixPlat,
		&flowNetwork,
		&priority,
		&producerHistory,
		&consumerHistory,
		&job,
		&stream,
		&commandGroup,
		&producer,
		&consumer,
		&writer,
		&writerClone,
		&integration,
		&mqttConnection,
		&credentials,
		&serialConnection,
		&ipConnection,
		&schedule,
		&tags,
	}

	for _, v := range models {
		err = db.AutoMigrate(v)
		if err != nil {
			panic("failed to AutoMigrate")
		}
	}

	var userCount int64 = 0
	db.Find(new(model.User)).Count(&userCount)
	if createDefaultUserIfNotExist && userCount == 0 {
		db.Create(&model.User{Name: defaultUser, Pass: password.CreatePassword(defaultPass, strength), Admin: true})
		if !production { //make a fake token for dev
			c := new(model.Client)
			c.Token = "fakeToken123"
			c.UserID = 1
			c.Name = "admin"
			db.Create(c)
		}

	}

	var platCount int64 = 0
	rp := new(model.RubixPlat)
	db.Find(rp).Count(&platCount)
	if createDefaultUserIfNotExist && platCount == 0 {
		rp.GlobalUuid = utils.MakeTopicUUID(model.CommonNaming.Rubix)
		db.Create(&rp)
	}
	busService := eventbus.NewService(eventbus.GetBus())
	gormDatabase = &GormDatabase{DB: db, Bus: busService}
	return &GormDatabase{DB: db, Bus: busService}, nil
}

func createDirectoryIfSqlite(dialect, connection string) {
	if dialect == "sqlite3" {
		if _, err := os.Stat(filepath.Dir(connection)); os.IsNotExist(err) {
			if err := mkdirAll(filepath.Dir(connection), 0777); err != nil {
				panic(err)
			}
		}
	}
}

// GormDatabase is a wrapper for the gorm framework.
type GormDatabase struct {
	DB            *gorm.DB
	Store         cachestore.Handler
	Bus           eventbus.BusService
	PluginManager *plugin.Manager
}

// Close closes the gorm database connection.
func (d *GormDatabase) Close() {
	fmt.Println("FIX THIS CLOSE DB dont work after upgrade of gorm")
	//d.Close() //TODO this is broken it calls itself
}
