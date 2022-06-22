package main

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/NubeIO/flow-framework/auth"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/database"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/history"
	"github.com/NubeIO/flow-framework/logger"
	"github.com/NubeIO/flow-framework/mqttclient"
	"github.com/NubeIO/flow-framework/router"
	"github.com/NubeIO/flow-framework/runner"
	"github.com/NubeIO/flow-framework/src/cachestore"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/NubeIO/flow-framework/src/jobs"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
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

func initHistory(db *database.GormDatabase, conf *config.Configuration) {
	h := new(history.History)
	h.DB = db
	if *conf.ProducerHistory.Cleaner.Enable {
		h.InitProducerHistoryCleaner(
			conf.ProducerHistory.Cleaner.Frequency,
			conf.ProducerHistory.Cleaner.DataPersistingHours)
	}
	if *conf.ProducerHistory.SyncInterval.Enable {
		h.InitProducerHistorySyncInterval(conf.ProducerHistory.SyncInterval.SyncPeriod)
	}
}

func main() {
	conf := config.CreateApp()
	logger.SetLogger(conf.LogLevel)
	logger.SetGinMode(conf.LogLevel)
	if err := os.MkdirAll(conf.GetAbsPluginDir(), 0755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(conf.GetAbsUploadedImagesDir(), 0755); err != nil {
		panic(err)
	}
	auth.CreateInternalTokenIfDoesNotExist()
	connection := path.Join(conf.GetAbsDataDir(), conf.Database.Connection)
	mqttBroker := "tcp://" + conf.MQTT.Address + ":" + strconv.Itoa(conf.MQTT.Port)
	_, err := mqttclient.InternalMQTT(mqttBroker)
	if err != nil {
		log.Errorln(err)
	}
	eventbus.Init()
	db, err := database.New(conf.Database.Dialect, connection, conf.Database.LogLevel)
	if err != nil {
		panic(err)
	}
	intHandler(db)
	defer db.Close()
	engine := router.Create(db, conf)
	eventbus.RegisterMQTTBus(false)
	initHistory(db, conf)
	runner.Run(engine, conf)
}
