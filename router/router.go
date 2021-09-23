package router

import (
	"fmt"
	"github.com/NubeDev/flow-framework/logger"
	"github.com/NubeDev/location"
	"github.com/gin-contrib/cors"
	"time"

	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/api/stream"
	"github.com/NubeDev/flow-framework/auth"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/error"
	"github.com/NubeDev/flow-framework/eventbus"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin"
	"github.com/gin-gonic/gin"
)

// Create creates the gin engine with all routes.
func Create(db *database.GormDatabase, vInfo *model.VersionInfo, conf *config.Configuration) (*gin.Engine, func()) {
	engine := gin.New()
	engine.Use(logger.GinMiddlewareLogger(), gin.Recovery(), error.Handler(), location.Default())
	engine.NoRoute(error.NotFound())
	eventBus := eventbus.NewService(eventbus.GetBus())
	streamHandler := stream.New(time.Duration(conf.Server.Stream.PingPeriodSeconds)*time.Second, 15*time.Second, conf.Server.Stream.AllowedOrigins, conf.Prod)
	//authentication := auth.Auth{DB: db}
	messageHandler := api.MessageAPI{Notifier: streamHandler, DB: db}
	healthHandler := api.HealthAPI{DB: db}
	clientHandler := api.ClientAPI{
		DB:            db,
		ImageDir:      conf.GetAbsUploadedImagesDir(),
		NotifyDeleted: streamHandler.NotifyDeletedClient,
	}
	applicationHandler := api.ApplicationAPI{
		DB:       db,
		ImageDir: conf.GetAbsUploadedImagesDir(),
	}
	userChangeNotifier := new(api.UserChangeNotifier)
	userHandler := api.UserAPI{DB: db, PasswordStrength: conf.PassStrength, UserChangeNotifier: userChangeNotifier}
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
	producerHandler := api.ProducerAPI{
		DB: db,
	}
	consumerHandler := api.ConsumersAPI{
		DB: db,
	}
	writerHandler := api.WriterAPI{
		DB: db,
	}
	writerCloneHandler := api.WriterCloneAPI{
		DB: db,
	}
	rubixPlatHandler := api.RubixPlatAPI{
		DB: db,
	}
	rubixCommandGroup := api.CommandGroupAPI{
		DB: db,
	}
	flowNetwork := api.FlowNetworksAPI{
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

	jobHandler.NewJobEngine()
	dbGroup.SyncTopics()
	//for the custom plugin endpoints you need to use the plugin token
	//http://0.0.0.0:1660/plugins/api/UUID/PLUGIN_TOKEN/echo
	pluginManager, err := plugin.NewManager(db, conf.GetAbsPluginDir(), engine.Group("/api/plugins/api/:uuid"), streamHandler)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.PluginManager = pluginManager
	pluginHandler := api.PluginAPI{
		Manager:  pluginManager,
		Notifier: streamHandler,
		DB:       db,
	}
	userChangeNotifier.OnUserDeleted(streamHandler.NotifyDeletedUser)
	userChangeNotifier.OnUserDeleted(pluginManager.RemoveUser)
	userChangeNotifier.OnUserAdded(pluginManager.InitializeForUserID)

	engine.GET("/api/system/ping", healthHandler.Health)
	engine.Static("/image", conf.GetAbsUploadedImagesDir())
	engine.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		for header, value := range conf.Server.ResponseHeaders {
			ctx.Header(header, value)
		}
	})
	engine.Use(cors.New(auth.CorsConfig(conf)))
	engine.OPTIONS("/*any")

	apiRoutes := engine.Group("/api")
	{
		apiRoutes.GET("/version", func(ctx *gin.Context) {
			ctx.JSON(200, vInfo)
		})

		//requireClientsGroupRoutes := apiRoutes.Group("", authentication.RequireClient())
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

			applicationRoutes := requireClientsGroupRoutes.Group("/applications")
			{
				applicationRoutes.GET("", applicationHandler.GetApplications)
				applicationRoutes.POST("", applicationHandler.CreateApplication)
				applicationRoutes.POST("/:id/image", applicationHandler.UploadApplicationImage)
				applicationRoutes.PUT("/:id", applicationHandler.UpdateApplication)
				applicationRoutes.DELETE("/:id", applicationHandler.DeleteApplication)

				tokenMessageRoutes := applicationRoutes.Group("/:id/messages")
				{
					tokenMessageRoutes.GET("", messageHandler.GetMessagesWithApplication)
					tokenMessageRoutes.DELETE("", messageHandler.DeleteMessageWithApplication)
				}
			}

			clientRoutes := requireClientsGroupRoutes.Group("/clients")
			{
				clientRoutes.GET("", clientHandler.GetClients)
				clientRoutes.POST("", clientHandler.CreateClient)
				clientRoutes.DELETE("/:id", clientHandler.DeleteClient)
				clientRoutes.PUT("/:id", clientHandler.UpdateClient)
			}

			messageRoutes := requireClientsGroupRoutes.Group("/messages")
			{
				messageRoutes.GET("", messageHandler.GetMessages)
				messageRoutes.DELETE("", messageHandler.DeleteMessages)
				messageRoutes.DELETE("/:id", messageHandler.DeleteMessage)
			}
		}

		//apiRoutes.Group("").Use(authentication.RequireApplicationToken()).POST("/messages", messageHandler.CreateMessage)
		apiRoutes.Group("").POST("/messages", messageHandler.CreateMessage)

		//apiRoutes.Use(authentication.RequireAdmin())

		userRoutes := apiRoutes.Group("/users")
		{
			userRoutes.GET("", userHandler.GetUsers)
			userRoutes.POST("", userHandler.CreateUser)
			userRoutes.PATCH("/current/password", userHandler.ChangePassword)
			userRoutes.GET("/current", userHandler.GetCurrentUser)
			userRoutes.DELETE("/:id", userHandler.DeleteUserByID)
			userRoutes.GET("/:id", userHandler.GetUserByID)
			userRoutes.PATCH("/:id", userHandler.UpdateUserByID)
		}

		databaseRoutes := apiRoutes.Group("/database")
		{
			//delete all networks, gateways, commandGroup, consumers, jobs and children.
			databaseRoutes.DELETE("/flows/drop", dbGroup.DropAllFlow)
			databaseWizard := databaseRoutes.Group("wizard")
			{
				databaseWizard.POST("/mappings/local/points", dbGroup.WizardLocalPointMapping)
				databaseWizard.POST("/mappings/remote/points", dbGroup.WizardRemotePointMapping)
				databaseWizard.POST("/mappings/existing/streams", dbGroup.Wizard2ndFlowNetwork)
				databaseWizard.POST("/networks/add", dbGroup.NetworkDevicePoint)
			}
		}

		wiresPlatRoutes := apiRoutes.Group("/wires/plat")
		{
			wiresPlatRoutes.GET("", rubixPlatHandler.GetRubixPlat)
			wiresPlatRoutes.PATCH("", rubixPlatHandler.UpdateRubixPlat)
		}

		historyProducerRoutes := apiRoutes.Group("/histories/producers")
		{
			historyProducerRoutes.GET("", historyHandler.GetProducerHistories)
			historyProducerRoutes.DELETE("/drop", historyHandler.DropProducerHistories)
			historyProducerRoutes.GET("/:uuid", historyHandler.GetProducerHistory)
			historyProducerRoutes.GET("/latest/:uuid", historyHandler.HistoryLatestByProducerUUID) //gets the newest
			historyProducerRoutes.GET("/all/:uuid", historyHandler.HistoriesAllByProducerUUID)
			historyProducerRoutes.POST("/bulk", historyHandler.CreateBulkProducerHistory)
			historyProducerRoutes.DELETE("/:uuid", historyHandler.DeleteProducerHistory)
		}

		flowNetworkRoutes := apiRoutes.Group("/flow/networks")
		{
			flowNetworkRoutes.GET("", flowNetwork.GetFlowNetworks)
			flowNetworkRoutes.POST("", flowNetwork.CreateFlowNetwork)
			flowNetworkRoutes.GET("/:uuid", flowNetwork.GetFlowNetwork)
			flowNetworkRoutes.PATCH("/:uuid", flowNetwork.UpdateFlowNetwork)
			flowNetworkRoutes.DELETE("/:uuid", flowNetwork.DeleteFlowNetwork)
			flowNetworkRoutes.GET("/one/args", flowNetwork.GetOneFlowNetworkByArgs)
			flowNetworkRoutes.DELETE("/drop", flowNetwork.DropFlowNetworks)
		}

		streamRoutes := apiRoutes.Group("/streams")
		{
			streamRoutes.GET("/", gatewayHandler.GetStreams)
			streamRoutes.POST("/", gatewayHandler.CreateStream)
			streamRoutes.GET("/:uuid", gatewayHandler.GetStream)
			streamRoutes.GET("/field/:uuid", gatewayHandler.GetStreamByField)
			streamRoutes.PATCH("/:uuid", gatewayHandler.UpdateStream)
			streamRoutes.DELETE("/:uuid", gatewayHandler.DeleteStream)
			streamRoutes.DELETE("/drop", gatewayHandler.DropStreams)
		}

		networkRoutes := apiRoutes.Group("/networks")
		{
			networkRoutes.GET("", networkHandler.GetNetworks)
			networkRoutes.POST("", networkHandler.CreateNetwork)
			networkRoutes.GET("/:uuid", networkHandler.GetNetwork)
			networkRoutes.GET("/plugin/:uuid", networkHandler.GetNetworkByPlugin)
			networkRoutes.GET("/plugin/all/:uuid", networkHandler.GetNetworksByPlugin)
			networkRoutes.PATCH("/:uuid", networkHandler.UpdateNetwork)
			networkRoutes.DELETE("/:uuid", networkHandler.DeleteNetwork)
			networkRoutes.DELETE("/drop", networkHandler.DropNetworks)
		}

		deviceRoutes := apiRoutes.Group("/devices")
		{
			deviceRoutes.GET("", deviceHandler.GetDevices)
			deviceRoutes.POST("", deviceHandler.CreateDevice)
			deviceRoutes.GET("/:uuid", deviceHandler.GetDevice)
			deviceRoutes.POST("/field/:uuid", deviceHandler.GetDeviceByField)
			deviceRoutes.PATCH("/field/:uuid", deviceHandler.UpdateDeviceByField)
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
			pointRoutes.GET("/field/:uuid", pointHandler.GetPointByField)
			pointRoutes.PATCH("/field/:uuid", pointHandler.UpdatePointByFieldAndType)
			pointRoutes.DELETE("/:uuid", pointHandler.DeletePoint)
			pointRoutes.DELETE("/drop", pointHandler.DropPoints)
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

			producerWriterCloneRoutes := producerRoutes.Group("/writers")
			{
				producerWriterCloneRoutes.GET("", writerCloneHandler.GetWriterClones)
				producerWriterCloneRoutes.POST("", writerCloneHandler.CreateWriterClone)
				producerWriterCloneRoutes.GET("/:uuid", writerCloneHandler.GetWriterClone)
				producerWriterCloneRoutes.PATCH("/:uuid", writerCloneHandler.UpdateWriterClone)
				producerWriterCloneRoutes.DELETE("/:uuid", writerCloneHandler.DeleteWriterClone)
				producerWriterCloneRoutes.DELETE("/drop", writerCloneHandler.DropWriterClone)
			}
		}

		consumerRoutes := apiRoutes.Group("/consumers")
		{
			consumerRoutes.GET("", consumerHandler.GetConsumers)
			consumerRoutes.POST("", consumerHandler.CreateConsumer)
			consumerRoutes.POST("/wizard", consumerHandler.AddConsumerWizard)
			consumerRoutes.GET("/:uuid", consumerHandler.GetConsumer)
			consumerRoutes.PATCH("/:uuid", consumerHandler.UpdateConsumer)
			consumerRoutes.DELETE("/:uuid", consumerHandler.DeleteConsumer)
			consumerRoutes.DELETE("/drop", consumerHandler.DropConsumers)

			consumerWriterRoutes := consumerRoutes.Group("/writers")
			{
				consumerWriterRoutes.GET("", writerHandler.GetWriters)
				consumerWriterRoutes.POST("wizard", writerHandler.CreateWriterWizard)
				consumerWriterRoutes.POST("", writerHandler.CreateWriter)
				consumerWriterRoutes.GET("/:uuid", writerHandler.GetWriter)
				consumerWriterRoutes.PATCH("/:uuid", writerHandler.UpdateWriter)
				consumerWriterRoutes.DELETE("/:uuid", writerHandler.DeleteWriter)
				consumerWriterRoutes.DELETE("/drop", writerHandler.DropWriters)
			}

		}

		//action's writers
		apiRoutes.POST("/writers/action/:uuid", writerHandler.WriterAction)
		apiRoutes.POST("/writers/action/bulk", writerHandler.WriterBulkAction)

		//action's writers clones
		apiRoutes.GET("/writers/clone/:uuid", writerCloneHandler.GetWriterClone)
		apiRoutes.PATCH("/writers/clone/:uuid", writerCloneHandler.UpdateWriterClone)

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
			integrationRoutes.GET("", integrationHandler.GetIntegrationsList)
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
			schRoutes.PATCH("/:uuid", schHandler.UpdateSchedule)
			schRoutes.DELETE("/:uuid", schHandler.DeleteSchedule)
			schRoutes.DELETE("/drop", schHandler.DropSchedules)
		}
		thingRoutes := apiRoutes.Group("/things")
		{
			thingRoutes.GET("/class", thingHandler.ThingClass)
			thingRoutes.GET("/types", thingHandler.ThingTypes)
			thingRoutes.GET("/writers/actions", thingHandler.WriterActions)
		}

		tagRoutes := apiRoutes.Group("/tags")
		{
			tagRoutes.GET("", tagHandler.GetTags)
			tagRoutes.POST("", tagHandler.CreateTag)
			tagRoutes.GET("/:tag", tagHandler.GetTag)
			tagRoutes.DELETE(":tag", tagHandler.DeleteTag)
		}

	}
	return engine, streamHandler.Close
}
