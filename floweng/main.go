package floweng

import (
	"flag"

	"github.com/NubeDev/flow-framework/floweng/server"
)

var (
	port = flag.String("port", "7071", "stream tools port")
)

var flowEngServer *server.Server

func FlowengStart() {

	flag.Parse()

	settings := server.NewSettings()
	flowEngServer = server.NewServer(settings)
	flowEngServer.SetRouter()
}
