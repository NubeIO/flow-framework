package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/migration"
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

	if err := migration.Upgrade(db); err != nil {
		panic("failed to AutoMigrate")
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
