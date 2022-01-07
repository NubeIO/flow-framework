package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/src/jobs"
	"os"
	"path"
	"time"

	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/database"
	"github.com/NubeIO/flow-framework/floweng"
	"github.com/NubeIO/flow-framework/logger"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/router"
	"github.com/NubeIO/flow-framework/runner"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var (
	Version   = "<version>"
	Commit    = "<commit>"
	BuildDate = "<build_date>"
)

func intHandler(db *database.GormDatabase) {
	dh := new(dbhandler.Handler)
	dh.DB = db
	dbhandler.Init(dh)

	s := new(cachestore.Handler)
	s.Store = cache.New(5*time.Minute, 10*time.Minute)
	cachestore.Init(s)

	j := new(jobs.Jobs)
	j.InitCron()
	if err := j.RefreshTokenJobAdd(); err != nil {
		panic(err)
	}
}

func main() {
	conf := config.CreateApp()
	logger.SetLogger(conf.LogLevel)
	logger.SetGinMode(conf.LogLevel)
	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	log.Info("Info Starting version:", vInfo.Version+"-"+vInfo.Commit+"@"+vInfo.BuildDate)
	if err := os.MkdirAll(conf.GetAbsPluginDir(), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(conf.GetAbsUploadedImagesDir(), 0755); err != nil {
		panic(err)
	}
	connection := path.Join(conf.GetAbsDataDir(), conf.Database.Connection)
	localBroker := "tcp://0.0.0.0:1883" //TODO add to config, this is meant to be an unsecure broker
	connected, err := mqttclient.InternalMQTT(localBroker)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("INIT INTERNAL MQTT CONNECTED", connected, "ERROR:", err)
	eventbus.Init()
	db, err := database.New(conf.Database.Dialect, connection, conf.Database.LogLevel)
	if err != nil {
		panic(err)
	}
	intHandler(db)
	floweng.EngStart(db)
	defer db.Close()
	engine, closeable := router.Create(db, vInfo, conf)
	defer closeable()
	eventbus.RegisterMQTTBus()
	runner.Run(engine, conf)
}
