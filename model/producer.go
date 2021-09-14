package model

//Producer a producer is a placeholder to register an object to enable consumers to
// A producer for example is a point, Something that makes data, and the subscriber would have a consumer to it, Like grafana reading and writing to it from edge to cloud or wires over rest(peer to peer)
type Producer struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	ProducerThingName     string            `json:"producer_thing_name"`
	ProducerThingUUID     string            `json:"producer_thing_uuid"` // this is the remote point UUID
	ProducerThingClass    string            `json:"producer_thing_class"`
	ProducerThingType     string            `json:"producer_thing_type"`
	StreamUUID            string            `json:"stream_uuid" gorm:"TYPE:string REFERENCES streams;not null;default:null"`
	CommonCurrentProducer                   //if the point for example is read only the writer uuid would be the point uuid, ie: itself, so in this case there is no writer or writer clone
	EnableHistory         bool              `json:"enable_history"`
	ProducerApplication   string            `json:"producer_application"`
	WriterClone           []WriterClone     `json:"writer_clones" gorm:"constraint:OnDelete:CASCADE;"`
	ProducerHistory       []ProducerHistory `json:"producer_histories" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}

//ProducerBody could be a local network, job or alarm and so on
type ProducerBody struct {
	CommonThingClass             //point, job
	CommonThingType              // for example temp, rssi, voltage
	FlowNetworkUUID  string      `json:"flow_network_uuid"`
	ProducerUUID     string      `json:"producer_uuid,omitempty"`
	StreamUUID       string      `json:"stream_uuid,omitempty"`
	Payload          interface{} `json:"payload"`
}
