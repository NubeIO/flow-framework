package floweng

import (
	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/floweng/core"
	"github.com/NubeDev/flow-framework/floweng/server"
)

var flowEngServer *server.Server

func EngStart(db *database.GormDatabase) {
	eventPool := make(chan core.BlockState, 256)
	flowEngServer = server.NewServer(db, eventPool)
	flowEngServer.LoadFromDB(db)
	flowEngServer.SetRouter()
	for i := 0; i < 10; i++ {
		go flowEngServer.RunRoutine()
	}
}
