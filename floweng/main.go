package floweng

import (
	"flag"
	"log"
	"net/http"

	"github.com/NubeDev/flow-framework/floweng/server"
)

var (
	port = flag.String("port", "7071", "stream tools port")
)

func main() {

	flag.Parse()

	settings := server.NewSettings()
	s := server.NewServer(settings)
	r := s.NewRouter()

	http.Handle("/", r)

	log.Println("serving on", *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Panicf(err.Error())
	}
}
