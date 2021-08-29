package router

import (
	"time"

	"github.com/NubeDev/flow-framework/api"
	"github.com/NubeDev/flow-framework/api/stream"
	"github.com/NubeDev/flow-framework/auth"
	"github.com/NubeDev/flow-framework/config"
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/docs"
	"github.com/NubeDev/flow-framework/error"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin"
	"github.com/NubeDev/location"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


// Create creates the gin engine with all routes.
func Create(db *database.GormDatabase, vInfo *model.VersionInfo, conf *config.Configuration) (*gin.Engine, func()) {
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery(), error.Handler(), location.Default())
	g.NoRoute(error.NotFound())

	streamHandler := stream.New(time.Duration(conf.Server.Stream.PingPeriodSeconds)*time.Second, 15*time.Second, conf.Server.Stream.AllowedOrigins, conf.Prod)
	authentication := auth.Auth{DB: db}
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
		DB: db,
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
	gatewayHandler := api.GatewayAPI{
		DB: db,
	}
	producerHandler := api.ProducerAPI{
		DB: db,
	}
	producerListHandler := api.ProducerListAPI{
		DB: db,
	}
	subscriptionHandler := api.SubscriptionsAPI{
		DB: db,
	}
	subscriptionListHandler := api.SubscriptionListAPI{
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
	jobHandler.NewJobEngine()
	dbGroup.SyncTopics()
	pluginManager, err := plugin.NewManager(db, conf.GetAbsPluginDir(), g.Group("/plugin/:uuid/custom/"), streamHandler);if err != nil {
		panic(err)
	}
	pluginHandler := api.PluginAPI{
		Manager:  pluginManager,
		Notifier: streamHandler,
		DB:       db,
	}

	userChangeNotifier.OnUserDeleted(streamHandler.NotifyDeletedUser)
	userChangeNotifier.OnUserDeleted(pluginManager.RemoveUser)
	userChangeNotifier.OnUserAdded(pluginManager.InitializeForUserID)

	g.GET("/ip", api.Hostname) //TODO remove
	g.GET("/health", healthHandler.Health)
	g.GET("/swagger", docs.Serve)
	g.Static("/image", conf.GetAbsUploadedImagesDir())
	g.GET("/docs", docs.UI)

	g.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		for header, value := range conf.Server.ResponseHeaders {
			ctx.Header(header, value)
		}
	})
	g.Use(cors.New(auth.CorsConfig(conf)))

	{
		g.GET("/plugin", authentication.RequireClient(), pluginHandler.GetPlugins)
		pluginRoute := g.Group("/plugin/", authentication.RequireClient())
		{
			pluginRoute.GET("/:uuid", pluginHandler.GetPlugin)
			pluginRoute.GET("/path/:path", pluginHandler.GetPluginByPath)
			pluginRoute.GET("/:uuid/config", pluginHandler.GetConfig)
			pluginRoute.POST("/:uuid/config", pluginHandler.UpdateConfig)
			pluginRoute.GET("/:uuid/display", pluginHandler.GetDisplay)
			pluginRoute.POST("/:uuid/enable", pluginHandler.EnablePlugin)
			pluginRoute.POST("/:uuid/disable", pluginHandler.DisablePlugin)
			pluginRoute.POST("/:uuid/network", pluginHandler.EnablePlugin)
		}
	}

	g.OPTIONS("/*any")

	// swagger:operation GET /version version getVersion
	//
	// Get version information.
	//
	// ---
	// produces: [application/json]
	// responses:
	//   200:
	//     description: Ok
	//     schema:
	//         $ref: "#/definitions/VersionInfo"
	g.GET("version", func(ctx *gin.Context) {
		ctx.JSON(200, vInfo)
	})

	g.Group("/").Use(authentication.RequireApplicationToken()).POST("/message", messageHandler.CreateMessage)

	clientAuth := g.Group("")
	{
		clientAuth.Use(authentication.RequireClient())
		app := clientAuth.Group("/application")
		{
			app.GET("", applicationHandler.GetApplications)
			app.POST("", applicationHandler.CreateApplication)
			app.POST("/:id/image", applicationHandler.UploadApplicationImage)
			app.PUT("/:id", applicationHandler.UpdateApplication)
			app.DELETE("/:id", applicationHandler.DeleteApplication)

			tokenMessage := app.Group("/:id/message")
			{
				tokenMessage.GET("", messageHandler.GetMessagesWithApplication)
				tokenMessage.DELETE("", messageHandler.DeleteMessageWithApplication)
			}
		}

		client := clientAuth.Group("/api/client")
		{
			client.GET("", clientHandler.GetClients)
			client.POST("", clientHandler.CreateClient)
			client.DELETE("/:id", clientHandler.DeleteClient)
			client.PUT("/:id", clientHandler.UpdateClient)
		}

		message := clientAuth.Group("/message")
		{
			message.GET("", messageHandler.GetMessages)
			message.DELETE("", messageHandler.DeleteMessages)
			message.DELETE("/:id", messageHandler.DeleteMessage)
		}

		clientAuth.GET("/stream", streamHandler.Handle)

	}
	authAdmin := g.Group("/api")
	{
		authAdmin.Use(authentication.RequireAdmin())
		authAdmin.GET("/users", userHandler.GetUsers)
		authAdmin.POST("/user", userHandler.CreateUser)
		authAdmin.PATCH("user/current/password", userHandler.ChangePassword)
		authAdmin.GET("/user/current", userHandler.GetCurrentUser)
		authAdmin.DELETE("/user/:id", userHandler.DeleteUserByID)
		authAdmin.GET("/user/:id", userHandler.GetUserByID)
		authAdmin.PATCH("/user/:id", userHandler.UpdateUserByID)
	}

	control := g.Group("api")
	{
		control.Use(authentication.RequireAdmin())
		control.GET("", api.Hostname)
		//delete all networks, gateways, commandGroup, subscriptions, jobs and children.
		control.DELETE("/database/flows/drop", dbGroup.DropAllFlow)
		control.POST("/database/wizard/mapping/local/point", dbGroup.WizardLocalPointMapping)
		control.POST("/database/wizard/mapping/remote/point", dbGroup.WizardRemotePointMapping)
		control.GET("/wires/plat", rubixPlatHandler.GetRubixPlat)
		control.PATCH("/wires/plat", rubixPlatHandler.UpdateRubixPlat)


		control.GET("/history/producers", historyHandler.GetProducerHistories)
		control.DELETE("/history/producers/drop", historyHandler.DropProducerHistories)
		control.GET("/history/producers/:uuid", historyHandler.GetProducerHistory)
		control.DELETE("/history/producers/:uuid", historyHandler.DeleteProducerHistory)

		control.GET("/networks", networkHandler.GetNetworks)
		control.DELETE("/networks/drop", networkHandler.DropNetworks)
		control.POST("/network", networkHandler.CreateNetwork)
		control.GET("/network/:uuid", networkHandler.GetNetwork)
		control.PATCH("/network/:uuid", networkHandler.UpdateNetwork)
		control.DELETE("/network/:uuid", networkHandler.DeleteNetwork)

		control.GET("/networks/flow", flowNetwork.GetFlowNetworks)
		control.DELETE("/networks/flow/drop", flowNetwork.DropFlowNetworks)
		control.POST("/network/flow", flowNetwork.CreateFlowNetwork)
		control.GET("/network/flow/:uuid", flowNetwork.GetFlowNetwork)
		control.PATCH("/network/flow/:uuid", flowNetwork.UpdateFlowNetwork)
		control.DELETE("/network/flow/:uuid", flowNetwork.DeleteFlowNetwork)

		control.GET("/devices", deviceHandler.GetDevices)
		control.DELETE("/devices/drop", deviceHandler.DropDevices)
		control.POST("/device", deviceHandler.CreateDevice)
		control.GET("/device/:uuid", deviceHandler.GetDevice)
		control.PATCH("/device/:uuid", deviceHandler.UpdateDevice)
		control.DELETE("/device/:uuid", deviceHandler.DeleteDevice)

		control.GET("/points", pointHandler.GetPoints)
		control.DELETE("/points/drop", pointHandler.DropPoints)
		control.POST("/point", pointHandler.CreatePoint)
		control.GET("/point/:uuid", pointHandler.GetPoint)
		control.PATCH("/point/:uuid", pointHandler.UpdatePoint)
		control.DELETE("/point/:uuid", pointHandler.DeletePoint)

		control.GET("/streams", gatewayHandler.GetStreamGateways)
		control.POST("/stream", gatewayHandler.CreateStreamGateway)
		control.GET("/stream/:uuid", gatewayHandler.GetStreamGateway)
		control.PATCH("/stream/:uuid", gatewayHandler.UpdateGateway)
		control.DELETE("/stream/:uuid", gatewayHandler.DeleteStreamGateway)

		control.GET("/commands", rubixCommandGroup.GetCommandGroups)
		control.POST("/command", rubixCommandGroup.CreateCommandGroup)
		control.GET("/command/:uuid", rubixCommandGroup.GetCommandGroup)
		control.PATCH("/command/:uuid", rubixCommandGroup.UpdateCommandGroup)
		control.DELETE("/command/:uuid", rubixCommandGroup.DeleteCommandGroup)


		control.GET("/producers", producerHandler.GetProducers)
		control.POST("/producer", producerHandler.CreateProducer)
		control.GET("/producer/:uuid", producerHandler.GetProducer)
		control.PATCH("/producer/:uuid", producerHandler.UpdateProducer)
		control.DELETE("/producer/:uuid", producerHandler.DeleteProducer)

		control.GET("/producers/list", producerListHandler.GetProducerLists)
		control.POST("/producer/list", producerListHandler.CreateProducerList)
		control.GET("/producer/list/:uuid", producerListHandler.GetProducerList)
		control.PATCH("/producer/list/:uuid", producerListHandler.UpdateProducerList)
		control.DELETE("/producer/list/:uuid", producerListHandler.DeleteProducerList)

		control.GET("/subscriptions", subscriptionHandler.GetSubscriptions)
		control.POST("/subscription", subscriptionHandler.CreateSubscription)
		control.GET("/subscription/:uuid", subscriptionHandler.GetSubscription)
		control.PATCH("/subscription/:uuid", subscriptionHandler.UpdateSubscription)
		control.DELETE("/subscription/:uuid", subscriptionHandler.DeleteSubscription)

		control.GET("/subscriptions/list", subscriptionListHandler.GetSubscriptionLists)
		control.POST("/subscription/list", subscriptionListHandler.CreateSubscriptionList)
		control.GET("/subscription/list/:uuid", subscriptionListHandler.GetSubscriptionList)
		control.PATCH("/subscription/list/:uuid", subscriptionListHandler.UpdateSubscriptionList)
		control.DELETE("/subscription/list/:uuid", subscriptionListHandler.DeleteSubscriptionList)

		//action's
		control.POST("/subscription/action/:uuid", subscriptionListHandler.SubscriptionAction)


		control.GET("/jobs", jobHandler.GetJobs)
		control.POST("/job", jobHandler.CreateJob)
		control.GET("/job/:uuid", jobHandler.GetJob)
		control.PATCH("/job/:uuid", jobHandler.UpdateJob)
		control.DELETE("/job/:uuid", jobHandler.DeleteJob)

	}

	return g, streamHandler.Close
}
