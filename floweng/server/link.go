package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/NubeDev/flow-framework/model"
	"github.com/gorilla/mux"
)

type LinkLedger struct {
	Source struct {
		Id int `json:"id"`
	} `json:"source"` // the source id
	Block struct {
		Id int `json:"id"`
	} `json:"block"` // the block id
	Id int `json:"id"` // link id
}

type ProtoLink struct {
	Id     int `json:"id"` // link id
	Source struct {
		Id int `json:"id"`
	} `json:"source"` // the source id
	Block struct {
		Id int `json:"id"`
	} `json:"block"` // the block id
}

func (s *Server) CreateLink(l ProtoLink) (*LinkLedger, error) {
	b, ok := s.blocks[l.Block.Id]
	if !ok {
		return nil, errors.New("could not find block")
	}

	sl, ok := s.sources[l.Source.Id]
	if !ok {
		return nil, errors.New("could not find source")
	}

	link := &LinkLedger{}
	link.Id = l.Id
	if link.Id == 0 {
		link.Id = s.GetNextID()
	}
	link.Source.Id = l.Source.Id
	link.Block.Id = l.Block.Id

	err := b.Block.SetSource(sl.Source)
	if err != nil {
		return nil, err
	}

	receiveChan := make(chan interface{})
	sl.Source.AddLink(b.Block, receiveChan)
	s.links[link.Id] = link

	s.websocketBroadcast(Update{Action: CREATE, Type: LINK, Data: wsLink{*link}})

	return link, nil
}

func (s *Server) DeleteLink(id int) error {
	link, ok := s.links[id]
	if !ok {
		return errors.New("could not find link")
	}

	block, ok := s.blocks[link.Block.Id]
	if !ok {
		return errors.New("could not find block")
	}

	source, ok := s.sources[link.Source.Id]
	if !ok {
		return errors.New("could not find source")
	}

	block.Block.SetSource(nil)
	source.Source.RemoveLink(block.Block)
	delete(s.links, id)
	EngDB.DeleteModel(id, model.Link{ID: id})

	s.websocketBroadcast(Update{Action: DELETE, Type: LINK, Data: wsLink{wsId{id}}})
	return nil
}

func (s *Server) listLinks() []LinkLedger {
	links := []LinkLedger{}
	for _, l := range s.links {
		links = append(links, *l)
	}
	return links
}

func (s *Server) LinkIndexHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	c := s.listLinks()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(c); err != nil {
		panic(err)
	}
}

func (s *Server) LinkCreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	var newLink ProtoLink
	json.Unmarshal(body, &newLink)

	s.Lock()
	defer s.Unlock()

	nl, err := s.CreateLink(newLink)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	EngDB.CreateModel(&model.Link{
		ID:       nl.Id,
		SourceID: nl.Source.Id,
		BlockID:  nl.Block.Id,
	})

	w.WriteHeader(http.StatusOK)
	writeJSON(w, nl)
}

func (s *Server) LinkDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromMux(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, err)
		return
	}

	s.Lock()
	defer s.Unlock()

	err = s.DeleteLink(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, Error{err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
