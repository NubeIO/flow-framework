package server

import (
	"log"
	"sync"

	"github.com/NubeDev/flow-framework/database"
	"github.com/NubeDev/flow-framework/floweng/core"
	"github.com/NubeDev/flow-framework/model"
)

const (
	// actions
	DELETE = "delete"
	RESET  = "reset"
	UPDATE = "update"
	CREATE = "create"
	INFO   = "info"
	// nodes
	BLOCK  = "block"
	GROUP  = "group"
	SOURCE = "source"
	// edges
	LINK       = "link"
	CONNECTION = "connection"
	// attributes
	CHILD      = "child"
	ROUTE      = "route"
	GROUPROUTE = "groupRoute"
	PARAM      = "param"
)

// The Server maintains a set of handlers that coordinate the creation of Nodes
type Server struct {
	// TODO these maps aren't strictly necessary, but save constantly performing depth first searches
	groups        map[int]*Group
	parents       map[int]int
	blocks        map[int]*BlockLedger      // these are the nodes
	connections   map[int]*ConnectionLedger // connections between the nodes
	sources       map[int]*SourceLedger
	links         map[int]*LinkLedger
	library       map[string]core.Spec
	sourceLibrary map[string]core.SourceSpec
	lastID        int
	addSocket     chan *socket
	delSocket     chan *socket
	broadcast     chan []byte
	emitChan      chan []byte
	sync.Mutex
}

var EngDB *database.GormDatabase

// NewServer starts a new Server. This object is immediately up and running.
func NewServer(db *database.GormDatabase) *Server {

	groups := make(map[int]*Group)
	groups[0] = &Group{
		Label:        "root",
		Id:           0,
		Children:     []int{},
		Parent:       nil,
		HiddenRoutes: make(map[string]struct{}),
	}

	blocks := make(map[int]*BlockLedger)
	connections := make(map[int]*ConnectionLedger)
	sources := make(map[int]*SourceLedger)
	links := make(map[int]*LinkLedger)
	library := core.GetLibrary()
	sourceLibrary := core.GetSources()
	parents := make(map[int]int)
	s := &Server{
		lastID:        0,
		parents:       parents,
		groups:        groups,
		blocks:        blocks,
		sourceLibrary: sourceLibrary,
		connections:   connections,
		library:       library,
		links:         links,
		sources:       sources,
		addSocket:     make(chan *socket),
		delSocket:     make(chan *socket),
		broadcast:     make(chan []byte),
		emitChan:      make(chan []byte),
	}

	// ws stuff
	log.Println("starting websocket handler")
	go s.websocketRouter()

	// db stuff
	EngDB = db

	var blockArr []*model.Block
	db.GetModelList(&blockArr)
	for _, block := range blockArr {
		_, err := s.CreateBlock(ProtoBlock{
			Label:    block.Label,
			Parent:   0,
			Type:     block.Type,
			Position: Position{block.PosX, block.PosY},
		})
		if err != nil {
			log.Println(err)
		}
	}
	var connArr []*model.Connection
	db.GetModelList(&connArr)
	for _, conn := range connArr {
		_, err := s.CreateConnection(ProtoConnection{
			Source: ConnectionNode{
				Id:    conn.SourceID,
				Route: conn.SourceRoute,
			},
			Target: ConnectionNode{
				Id:    conn.TargetID,
				Route: conn.TargetRoute,
			},
		})
		if err != nil {
			log.Println(err)
		}
	}
	blockInputs, _ := db.GetBlockStaticInputs()
	for _, r := range blockInputs {
		err := s.ModifyBlockRoute(r.BlockID, r.RouteIndex, &core.InputValue{Data: r.Value})
		if err != nil {
			log.Println(err)
		}
	}
	// TODO: handle IDs changing after deleting nodes
	return s
}

// GetNextID returns the next ID to be used for a new group or a new block
func (s *Server) GetNextID() int {
	s.lastID++
	return s.lastID
}
