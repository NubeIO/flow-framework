package floweng

import (
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/floweng/server"
)

var flowEngServer *server.Server

func FlowengStart(db *database.GormDatabase) {
	flowEngServer = server.NewServer(db)
	flowEngServer.LoadFromDB(db)
	flowEngServer.SetRouter()
}
