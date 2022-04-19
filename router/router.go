package router

import (
	"fmt"
	"github.com/NubeDev/location"
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/flow-framework/auth"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/database"
	"github.com/NubeIO/flow-framework/error"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/logger"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/networking/networking"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Create(db *database.GormDatabase, conf *config.Configuration) *gin.Engine {
	engine := gin.New()
	engine.Use(logger.GinMiddlewareLogger(), gin.Recovery(), error.Handler(), location.Default())
	engine.NoRoute(error.NotFound())
	eventBus := eventbus.NewService(eventbus.GetBus())
	proxyHandler := api.Proxy{DB: db}
	healthHandler := api.HealthAPI{DB: db}
	localStorageFlowNetworkHandler := api.LocalStorageFlowNetworkAPI{
		DB: db,
	}
	networkHandler := api.NetworksAPI{
		DB:  db,
		Bus: eventBus,
	}
	deviceHandler := api.DeviceAPI{
		DB: db,
	}
	pointHandler := api.PointAPI{
		DB: db,
	}
	historyHandler := api.HistoriesAPI{
		DB: db,
	}
	jobHandler := api.JobAPI{
		DB: db,
	}
	gatewayHandler := api.StreamAPI{
		DB: db,
	}
	streamCloneHandler := api.StreamCloneAPI{
		DB: db,
	}
	producerHandler := api.ProducerAPI{
		DB: db,
	}
	consumerHandler := api.ConsumersAPI{
		DB: db,
	}
	writerCloneHandler := api.WriterCloneAPI{
		DB: db,
	}
	rubixCommandGroup := api.CommandGroupAPI{
		DB: db,
	}
	flowNetwork := api.FlowNetworksAPI{
		DB: db,
	}
	mapping := api.MappingAPI{
		DB: db,
	}
	flowNetworkCloneHandler := api.FlowNetworkClonesAPI{
		DB: db,
	}
	dbGroup := api.DatabaseAPI{
		DB: db,
	}
	integrationHandler := api.IntegrationAPI{
		DB: db,
	}
	mqttHandler := api.MqttConnectionAPI{
		DB: db,
	}
	schHandler := api.ScheduleAPI{
		DB: db,
	}
	thingHandler := api.ThingAPI{}

	tagHandler := api.TagAPI{
		DB: db,
	}
	nets := networking.NewNets()
	deviceInfoHandler := api.DeviceInfoAPI{
		DB:   db,
		Nets: nets,
	}
	writerHandler := api.WriterAPI{
		DB: db,
	}
	syncFlowNetworkHandler := api.SyncFlowNetworkAPI{
		DB: db,
	}
	syncStreamHandler := api.SyncStreamAPI{
		DB: db,
	}
	syncWriterHandler := api.SyncWriterAPI{
		DB: db,
	}
	dbGroup.SyncTopics()
	// for the custom plugin endpoints you need to use the plugin token
	// http://0.0.0.0:1660/plugins/api/UUID/PLUGIN_TOKEN/echo
	pluginManager, err := plugin.NewManager(db, conf.GetAbsPluginDir(), engine.Group("/api/plugins/api"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.PluginManager = pluginManager
	pluginHandler := api.PluginAPI{
		Manager: pluginManager,
		DB:      db,
	}

	engine.GET("/api/system/ping", healthHandler.Health)
	engine.Static("/image", conf.GetAbsUploadedImagesDir())
	engine.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json") // if you comment it out it will detected as text on proxy-handlers
		for header, value := range conf.Server.ResponseHeaders {
			ctx.Header(header, value)
		}
	})

	engine.Use(cors.New(auth.CorsConfig(conf)))
	engine.OPTIONS("/*any")

	apiRoutes := engine.Group("/api")
	{
		fnProxy := apiRoutes.Group("/fn")
		{
			fnProxy.GET("/*any", proxyHandler.GetProxy(true))
			fnProxy.POST("/*any", proxyHandler.PostProxy(true))
			fnProxy.PUT("/*any", proxyHandler.PutProxy(true))
			fnProxy.PATCH("/*any", proxyHandler.PatchProxy(true))
			fnProxy.DELETE("/*any", proxyHandler.DeleteProxy(true))
		}

		fncProxy := apiRoutes.Group("/fnc")
		{
			fncProxy.GET("/*any", proxyHandler.GetProxy(false))
			fncProxy.POST("/*any", proxyHandler.PostProxy(false))
			fncProxy.PUT("/*any", proxyHandler.PutProxy(false))
			fncProxy.PATCH("/*any", proxyHandler.PatchProxy(false))
			fncProxy.DELETE("/*any", proxyHandler.DeleteProxy(false))
		}

		requireClientsGroupRoutes := apiRoutes.Group("")
		{
			plugins := requireClientsGroupRoutes.Group("/plugins")
			{
				plugins.GET("", pluginHandler.GetPlugins)
				plugins.GET("/:uuid", pluginHandler.GetPlugin)
				plugins.GET("/config/:uuid", pluginHandler.GetConfig)
				plugins.POST("/config/:uuid", pluginHandler.UpdateConfig)
				plugins.GET("/display/:uuid", pluginHandler.GetDisplay)
				plugins.POST("/enable/:uuid", pluginHandler.EnablePluginByUUID)
				plugins.POST("/restart/:uuid", pluginHandler.RestartPlugin)
				plugins.GET("/path/:path", pluginHandler.GetPluginByPath)
			}
		}

		databaseRoutes := apiRoutes.Group("/database")
		{
			// delete all networks, gateways, commandGroup, consumers, jobs and children
			databaseRoutes.DELETE("/flows/drop", dbGroup.DropAllFlow)
			databaseWizard := databaseRoutes.Group("wizard")
			{
				databaseWizard.POST("/mappings/p2p/points", dbGroup.WizardP2PMapping)
				databaseWizard.POST("/mappings/master_slave/points", dbGroup.WizardMasterSlavePointMapping)
				databaseWizard.POST("/mapping/master_slave/points/consumer/:global_uuid", dbGroup.WizardMasterSlavePointMappingOnConsumerSideByProducerSide) //supplementary API for remote_mapping
				databaseWizard.POST("/mapping/p2p/points/consumer/:global_uuid", dbGroup.WizardP2PMappingOnConsumerSideByProducerSide)                       //supplementary API for remote_mapping
			}
		}

		localStorageFlowNetworkRoutes := apiRoutes.Group("/localstorage_flow_network")
		{
			localStorageFlowNetworkRoutes.GET("", localStorageFlowNetworkHandler.GetLocalStorageFlowNetwork)
			localStorageFlowNetworkRoutes.PATCH("", localStorageFlowNetworkHandler.UpdateLocalStorageFlowNetwork)
			localStorageFlowNetworkRoutes.GET("/refresh_flow_token", localStorageFlowNetworkHandler.RefreshLocalStorageFlowToken)
		}

		historyProducerRoutes := apiRoutes.Group("/histories/producers")
		{
			historyProducerRoutes.GET("", historyHandler.GetProducerHistories)
			historyProducerRoutes.GET("/:producer_uuid", historyHandler.GetProducerHistoriesByProducerUUID)
			historyProducerRoutes.GET("/name/:name/one", historyHandler.GetLatestProducerHistoryByProducerName)
			historyProducerRoutes.GET("/:producer_uuid/one", historyHandler.GetLatestProducerHistoryByProducerUUID)
			historyProducerRoutes.GET("/points", historyHandler.GetProducerHistoriesPoints)
			historyProducerRoutes.POST("", historyHandler.CreateProducerHistory)
			historyProducerRoutes.DELETE("/all", historyHandler.DeleteAllProducerHistories)
			historyProducerRoutes.DELETE("/:producer_uuid", historyHandler.DeleteProducerHistoriesByProducerUUID)
		}

		flowNetworkRoutes := apiRoutes.Group("/flow_networks")
		{
			flowNetworkRoutes.GET("", flowNetwork.GetFlowNetworks)
			flowNetworkRoutes.POST("", flowNetwork.CreateFlowNetwork)
			flowNetworkRoutes.GET("/:uuid", flowNetwork.GetFlowNetwork)
			flowNetworkRoutes.PATCH("/:uuid", flowNetwork.UpdateFlowNetwork)
			flowNetworkRoutes.DELETE("/:uuid", flowNetwork.DeleteFlowNetwork)
			flowNetworkRoutes.GET("/one/args", flowNetwork.GetOneFlowNetworkByArgs)
			flowNetworkRoutes.DELETE("/drop", flowNetwork.DropFlowNetworks)
			flowNetworkRoutes.GET("/refresh_connections", flowNetwork.RefreshFlowNetworksConnections)
		}

		flowNetworkCloneRoutes := apiRoutes.Group("/flow_network_clones")
		{
			flowNetworkCloneRoutes.GET("", flowNetworkCloneHandler.GetFlowNetworkClones)
			flowNetworkCloneRoutes.GET("/:uuid", flowNetworkCloneHandler.GetFlowNetworkClone)
			flowNetworkCloneRoutes.DELETE("/:uuid", flowNetworkCloneHandler.DeleteFlowNetworkClone)
			flowNetworkCloneRoutes.GET("/one/args", flowNetworkCloneHandler.GetOneFlowNetworkCloneByArgs)
			flowNetworkCloneRoutes.GET("/refresh_connections", flowNetworkCloneHandler.RefreshFlowNetworkClonesConnections)
		}

		streamRoutes := apiRoutes.Group("/streams")
		{
			streamRoutes.GET("/", gatewayHandler.GetStreams)
			streamRoutes.POST("/", gatewayHandler.CreateStream)
			streamRoutes.GET("/:uuid", gatewayHandler.GetStream)
			streamRoutes.PATCH("/:uuid", gatewayHandler.UpdateStream)
			streamRoutes.DELETE("/:uuid", gatewayHandler.DeleteStream)
			streamRoutes.DELETE("/drop", gatewayHandler.DropStreams)
		}

		mappingRoutes := apiRoutes.Group("/mapping")
		{
			mappingRoutes.POST("/points", mapping.CreatePointMapping)
		}

		streamCloneRoutes := apiRoutes.Group("/stream_clones")
		{
			streamCloneRoutes.GET("", streamCloneHandler.GetStreamClones)
			streamCloneRoutes.GET("/:uuid", streamCloneHandler.GetStreamClone)
			streamCloneRoutes.DELETE("/:uuid", streamCloneHandler.DeleteStreamClone)
		}

		networkRoutes := apiRoutes.Group("/networks")
		{
			networkRoutes.GET("", networkHandler.GetNetworks)
			networkRoutes.POST("", networkHandler.CreateNetwork)
			networkRoutes.GET("/:uuid", networkHandler.GetNetwork)
			networkRoutes.GET("/plugin/:name", networkHandler.GetNetworkByPluginName)
			networkRoutes.GET("/plugin/all/:name", networkHandler.GetNetworksByPluginName)
			networkRoutes.GET("/name/:name", networkHandler.GetNetworkByName)
			networkRoutes.GET("/name/all/:name", networkHandler.GetNetworksByName)
			networkRoutes.PATCH("/:uuid", networkHandler.UpdateNetwork)
			networkRoutes.DELETE("/:uuid", networkHandler.DeleteNetwork)
			networkRoutes.DELETE("/drop", networkHandler.DropNetworks)
		}

		deviceRoutes := apiRoutes.Group("/devices")
		{
			deviceRoutes.GET("", deviceHandler.GetDevices)
			deviceRoutes.POST("", deviceHandler.CreateDevice)
			deviceRoutes.GET("/:uuid", deviceHandler.GetDevice)
			deviceRoutes.GET("/one/args", deviceHandler.GetOneDeviceByArgs)
			deviceRoutes.PATCH("/:uuid", deviceHandler.UpdateDevice)
			deviceRoutes.DELETE("/:uuid", deviceHandler.DeleteDevice)
			deviceRoutes.DELETE("/drop", deviceHandler.DropDevices)
		}

		pointRoutes := apiRoutes.Group("/points")
		{
			pointRoutes.GET("", pointHandler.GetPoints)
			pointRoutes.POST("", pointHandler.CreatePoint)
			pointRoutes.GET("/:uuid", pointHandler.GetPoint)
			pointRoutes.PATCH("/:uuid", pointHandler.UpdatePoint)
			pointRoutes.PATCH("/write/:uuid", pointHandler.PointWrite)
			pointRoutes.GET("/one/args", pointHandler.GetOnePointByArgs)
			pointRoutes.DELETE("/:uuid", pointHandler.DeletePoint)
			pointRoutes.DELETE("/drop", pointHandler.DropPoints)
			pointRoutes.GET("/name", pointHandler.GetPointByName)
			pointRoutes.PATCH("/name", pointHandler.PointWriteByName)
		}

		commandRoutes := apiRoutes.Group("/commands")
		{
			commandRoutes.GET("", rubixCommandGroup.GetCommandGroups)
			commandRoutes.POST("", rubixCommandGroup.CreateCommandGroup)
			commandRoutes.GET("/:uuid", rubixCommandGroup.GetCommandGroup)
			commandRoutes.PATCH("/:uuid", rubixCommandGroup.UpdateCommandGroup)
			commandRoutes.DELETE("/:uuid", rubixCommandGroup.DeleteCommandGroup)
		}

		producerRoutes := apiRoutes.Group("/producers")
		{
			producerRoutes.GET("", producerHandler.GetProducers)
			producerRoutes.POST("", producerHandler.CreateProducer)
			producerRoutes.GET("/:uuid", producerHandler.GetProducer)
			producerRoutes.PATCH("/:uuid", producerHandler.UpdateProducer)
			producerRoutes.DELETE("/:uuid", producerHandler.DeleteProducer)
			producerRoutes.DELETE("/drop", producerHandler.DropProducers)
			producerRoutes.GET("/one/args", producerHandler.GetOneProducerByArgs)

			producerWriterCloneRoutes := producerRoutes.Group("/writer_clones")
			{
				producerWriterCloneRoutes.GET("", writerCloneHandler.GetWriterClones)
				producerWriterCloneRoutes.POST("", writerCloneHandler.CreateWriterClone)
				producerWriterCloneRoutes.GET("/:uuid", writerCloneHandler.GetWriterClone)
				producerWriterCloneRoutes.DELETE("/:uuid", writerCloneHandler.DeleteWriterClone)
				producerWriterCloneRoutes.DELETE("/drop", writerCloneHandler.DropWriterClone)
			}
		}

		consumerRoutes := apiRoutes.Group("/consumers")
		{
			consumerRoutes.GET("", consumerHandler.GetConsumers)
			consumerRoutes.POST("", consumerHandler.CreateConsumer)
			consumerRoutes.GET("/:uuid", consumerHandler.GetConsumer)
			consumerRoutes.PATCH("/:uuid", consumerHandler.UpdateConsumer)
			consumerRoutes.DELETE("/:uuid", consumerHandler.DeleteConsumer)
			consumerRoutes.DELETE("/drop", consumerHandler.DropConsumers)

			consumerWriterRoutes := consumerRoutes.Group("/writers")
			{
				consumerWriterRoutes.GET("", writerHandler.GetWriters)
				consumerWriterRoutes.POST("", writerHandler.CreateWriter)
				consumerWriterRoutes.GET("/:uuid", writerHandler.GetWriter)
				consumerWriterRoutes.PATCH("/:uuid", writerHandler.UpdateWriter)
				consumerWriterRoutes.DELETE("/:uuid", writerHandler.DeleteWriter)
				consumerWriterRoutes.DELETE("/drop", writerHandler.DropWriters)
			}
		}

		jobRoutes := apiRoutes.Group("/jobs")
		{
			jobRoutes.GET("", jobHandler.GetJobs)
			jobRoutes.POST("", jobHandler.CreateJob)
			jobRoutes.GET("/:uuid", jobHandler.GetJob)
			jobRoutes.PATCH("/:uuid", jobHandler.UpdateJob)
			jobRoutes.DELETE("/:uuid", jobHandler.DeleteJob)
		}

		integrationRoutes := apiRoutes.Group("/integrations")
		{
			integrationRoutes.GET("", integrationHandler.GetIntegrations)
			integrationRoutes.POST("", integrationHandler.CreateIntegration)
			integrationRoutes.GET("/:uuid", integrationHandler.GetIntegration)
			integrationRoutes.PATCH("/:uuid", integrationHandler.UpdateIntegration)
			integrationRoutes.DELETE("/:uuid", integrationHandler.DeleteIntegration)
			integrationRoutes.DELETE("/drop", integrationHandler.DropIntegrationsList)
		}

		mqttClientRoutes := apiRoutes.Group("/mqtt/clients")
		{
			mqttClientRoutes.GET("", mqttHandler.GetMqttConnectionsList)
			mqttClientRoutes.POST("", mqttHandler.CreateMqttConnection)
			mqttClientRoutes.GET("/:uuid", mqttHandler.GetMqttConnection)
			mqttClientRoutes.PATCH("/:uuid", mqttHandler.UpdateMqttConnection)
			mqttClientRoutes.DELETE("/:uuid", mqttHandler.DeleteMqttConnection)
			mqttClientRoutes.DELETE("/drop", mqttHandler.DropMqttConnectionsList)
		}

		schRoutes := apiRoutes.Group("/schedules")
		{
			schRoutes.GET("", schHandler.GetSchedules)
			schRoutes.POST("", schHandler.CreateSchedule)
			schRoutes.GET("/:uuid", schHandler.GetSchedule)
			schRoutes.GET("/one/args", schHandler.GetOneScheduleByArgs)
			schRoutes.PATCH("/:uuid", schHandler.UpdateSchedule)
			schRoutes.PATCH("/write/:uuid", schHandler.ScheduleWrite)
			schRoutes.DELETE("/:uuid", schHandler.DeleteSchedule)
			schRoutes.DELETE("/drop", schHandler.DropSchedules)
		}

		thingRoutes := apiRoutes.Group("/things")
		{
			thingRoutes.GET("/class", thingHandler.ThingClass)
			thingRoutes.GET("/writers/actions", thingHandler.WriterActions)
			thingRoutes.GET("/units", thingHandler.ThingUnits)
		}

		tagRoutes := apiRoutes.Group("/tags")
		{
			tagRoutes.GET("", tagHandler.GetTags)
			tagRoutes.POST("", tagHandler.CreateTag)
			tagRoutes.GET("/:tag", tagHandler.GetTag)
			tagRoutes.DELETE(":tag", tagHandler.DeleteTag)
		}

		deviceInfoRoutes := apiRoutes.Group("/system")
		{
			deviceInfoRoutes.GET("/device_info", deviceInfoHandler.GetDeviceInfo)
			deviceInfoRoutes.GET("/host/info", deviceInfoHandler.GetOSDetails)
			deviceInfoRoutes.GET("/host/timezone/list", deviceInfoHandler.GetTZoneList)
			deviceInfoRoutes.GET("/time", deviceInfoHandler.GetSystemTime)
			deviceInfoRoutes.GET("/ip/internal", deviceInfoHandler.GetNetworks)
			deviceInfoRoutes.GET("/ip/external", deviceInfoHandler.GetExternalIP)
			deviceInfoRoutes.GET("/ip/interfaces", deviceInfoHandler.GetInterfacesNames)
			deviceInfoRoutes.GET("/firewall/status", deviceInfoHandler.FirewallStatus)
			deviceInfoRoutes.POST("/discover", deviceInfoHandler.RubixNetworkPing)
		}

		apiRoutes.POST("/writers/action/:uuid", writerHandler.WriterAction)
		apiRoutes.POST("/writers/action/bulk", writerHandler.WriterBulkAction)

		syncRoutes := apiRoutes.Group("/sync")
		{
			syncRoutes.POST("/flow_network", syncFlowNetworkHandler.SyncFlowNetwork)
			syncRoutes.POST("/stream", syncStreamHandler.SyncStream)
			syncRoutes.POST("/writer", syncWriterHandler.SyncWriter)
			syncRoutes.POST("/cov/:writer_uuid", syncWriterHandler.SyncCOV) // clone ---> source side
			syncRoutes.POST("/writer/write/:source_uuid", syncWriterHandler.SyncWriterWriteAction)
			syncRoutes.GET("/writer/read/:source_uuid", syncWriterHandler.SyncWriterReadAction)
		}
	}
	return engine
}
