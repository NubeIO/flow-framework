package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/NubeIO/flow-framework/logger"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var mkdirAll = os.MkdirAll
var gormDatabase *GormDatabase

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection, logLevel string) (*GormDatabase, error) {
	createDirectoryIfSqlite(dialect, connection)
	_connection := fmt.Sprintf("%s?_foreign_keys=on", connection)
	db, err := gorm.Open(sqlite.Open(_connection), &gorm.Config{
		Logger: logger.New().SetLogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		panic("failed to connect database")
	}
	var localStorageFlowNetwork []model.LocalStorageFlowNetwork
	var alerts []model.Alert
	var message []model.Message
	var pluginConf []model.PluginConf
	var network []model.Network
	var device []model.Device
	var point []model.Point
	var priority []model.Priority
	var producerHistory []model.ProducerHistory
	var consumerHistory []model.ConsumerHistory
	var flowNetwork []model.FlowNetwork
	var flowNetworkClone []model.FlowNetworkClone
	var job []model.Job
	var stream []model.Stream
	var streamClone []model.StreamClone
	var commandGroup []model.CommandGroup
	var producer []model.Producer
	var consumer []model.Consumer
	var writer []model.Writer
	var writerClone []model.WriterClone
	var integration []model.Integration
	var mqttConnection []model.MqttConnection
	var schedule []model.Schedule
	var tags []model.Tag
	var blocks []model.Block
	var connections []model.Connection
	var blockRoutes []model.BlockStaticRoute
	var blockRouteValueNumbers []model.BlockRouteValueNumber
	var blockRouteValueString []model.BlockRouteValueString
	var blockRouteValueBool []model.BlockRouteValueBool
	var sourceParams []model.SourceParameter
	var links []model.Link
	var history []model.History
	var historyLog []model.HistoryLog
	var historyInfluxLog []model.HistoryInfluxLog
	var models = []interface{}{
		&localStorageFlowNetwork,
		&alerts,
		&message,
		&pluginConf,
		&network,
		&device,
		&point,
		&flowNetwork,
		&flowNetworkClone,
		&priority,
		&producerHistory,
		&consumerHistory,
		&job,
		&stream,
		&streamClone,
		&commandGroup,
		&producer,
		&consumer,
		&writer,
		&writerClone,
		&integration,
		&mqttConnection,
		&schedule,
		&tags,
		&blocks,
		&connections,
		&blockRoutes,
		&blockRouteValueNumbers,
		&blockRouteValueString,
		&blockRouteValueBool,
		&sourceParams,
		&links,
		&history,
		&historyLog,
		&historyInfluxLog,
	}

	for _, v := range models {
		err = db.AutoMigrate(v)
		if err != nil {
			panic("failed to AutoMigrate")
		}
	}

	var lsFlowNetworkCount int64 = 0
	conf := config.Get()
	lsFlowNetwork := model.LocalStorageFlowNetwork{
		FlowIP:    conf.Server.ListenAddr,
		FlowPort:  conf.Server.RSPort,
		FlowHTTPS: utils.NewFalse(),
	}
	db.Find(&model.LocalStorageFlowNetwork{}).Count(&lsFlowNetworkCount)
	if lsFlowNetworkCount == 0 {
		db.Create(&lsFlowNetwork)
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
	sqlDB, err := d.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
