package model

/*
Producer a producer is a placeholder to register an object to enable consumers to a producer for example is a point,
something that makes data, and the subscriber would have a consumer to it, like grafana reading and writing to it
from edge to cloud or wires over rest(peer to peer)
*/
type Producer struct {
	CommonUUID
	CommonNameUnique
	CommonDescription
	CommonEnable
	CommonSyncUUID
	ProducerThingName       string             `json:"producer_thing_name,omitempty"`  //e.g. point.name, user will understand what name it is
	ProducerThingUUID       string             `json:"producer_thing_uuid,omitempty"`  //e.g. point.uuid
	ProducerThingClass      string             `json:"producer_thing_class,omitempty"` //e.g. point.thing_class, i.e. point, job etc.
	ProducerThingType       string             `json:"producer_thing_type,omitempty"`  //e.g. point.thing_type, i.e. temp, rssi, voltage etc.
	ProducerApplication     string             `json:"producer_application,omitempty"`
	CommonCurrentWriterUUID                    //if the point for example is read only the writer.uuid would be the point.uuid, i.e.: itself, so in this case there is no writer or writer clone
	EnableHistory           *bool              `json:"enable_history,omitempty"`
	StreamUUID              string             `json:"stream_uuid,omitempty" gorm:"TYPE:string REFERENCES streams;not null"`
	WriterClones            []*WriterClone     `json:"writer_clones,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	ProducerHistories       []*ProducerHistory `json:"producer_histories,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Tags                    []*Tag             `json:"tags,omitempty" gorm:"many2many:producers_tags;constraint:OnDelete:CASCADE"`
	CommonCreated
}

//ProducerBody could be a local network, job or alarm and so on
type ProducerBody struct {
	CommonThingClass //point, job
	CommonThingRef
	CommonThingType             //for example temp, rssi, voltage
	FlowNetworkUUID string      `json:"flow_network_uuid,omitempty"`
	ProducerUUID    string      `json:"producer_uuid,omitempty"`
	StreamUUID      string      `json:"stream_uuid,omitempty"`
	Payload         interface{} `json:"payload"`
}
