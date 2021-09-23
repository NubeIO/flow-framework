package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/cachestore"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/dbhandler"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/floweng"
	"github.com/NubeDev/flow-framework/logger"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/mqttclient"
	"github.com/NubeDev/flow-framework/router"
	"github.com/NubeDev/flow-framework/runner"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var (
	Version   = "<version>"
	Commit    = "<commit>"
	BuildDate = "<build_date>"
)

func intHandler(db *database.GormDatabase) {
	//db access
	dh := new(dbhandler.Handler)
	dh.DB = db
	dbhandler.Init(dh)
	//store access
	s := new(cachestore.Handler)
	s.Store = cache.New(5*time.Minute, 10*time.Minute)
	cachestore.Init(s)

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
	eventbus.RegisterMQTTBus()

	db, err := database.New(conf.Database.Dialect, connection, conf.DefaultUser.Name, conf.DefaultUser.Pass,
		conf.PassStrength, conf.Database.LogLevel, true, conf.Prod)
	if err != nil {
		panic(err)
	}
	p := conf.GetAbsConfigDir()
	thingType := fmt.Sprintf("%s/tags.json", p)
	plan, _ := ioutil.ReadFile(thingType)
	err = json.Unmarshal(plan, &api.ThingTypes)
	if err != nil {
		fmt.Println(err)
	}
	intHandler(db)
	floweng.FlowengStart(db)
	defer db.Close()
	engine, closeable := router.Create(db, vInfo, conf)
	defer closeable()
	runner.Run(engine, conf)

}
