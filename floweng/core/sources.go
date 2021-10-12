package core

func GetSources() map[string]SourceSpec {
	sources := []SourceSpec{
		// KeyValueStore(),
		// ValueStore(),
		// PriorityQueueStore(),
		// ListStore(),
		// WebsocketClient(),
		StdinInterface(),
	}

	library := make(map[string]SourceSpec)

	for _, s := range sources {
		library[s.Name] = s
	}

	return library
}

func (sc *SourceCommon) AddLink(block *Block, l chan interface{}) {
	sc.Links[block] = l
}

func (sc *SourceCommon) RemoveLink(block *Block, l chan interface{}) {
	delete(sc.Links, block)
}
