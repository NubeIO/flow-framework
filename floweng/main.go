package floweng

import (
	"github.com/NubeDev/flow-framework/floweng/server"
)

var flowEngServer *server.Server

func FlowengStart() {
	flowEngServer = server.NewServer()
	flowEngServer.SetRouter()
}
