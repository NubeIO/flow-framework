package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/migration"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/utils/boolean"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"os"
	"path/filepath"
	"time"

	"github.com/NubeIO/flow-framework/logger"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	username = "admin"
	password = "N00BWires"
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
	// db.Exec("PRAGMA journal_mode=WAL;")

	if err = migration.AutoMigrate(db); err != nil {
		panic("failed to AutoMigrate")
	}

	var lsFlowNetworkCount int64 = 0
	conf := config.Get()
	lsFlowNetwork := model.LocalStorageFlowNetwork{
		FlowIP:    conf.Server.ListenAddr,
		FlowPort:  conf.Server.RSPort,
		FlowHTTPS: boolean.NewFalse(),
	}
	db.Find(&model.LocalStorageFlowNetwork{}).Count(&lsFlowNetworkCount)
	if lsFlowNetworkCount == 0 {
		db.Create(&lsFlowNetwork)
	}

	user_, _ := user.GetUser()
	if user_ == nil {
		_, _ = user.CreateUser(&user.User{Username: username, Password: password})
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
