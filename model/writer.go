package model

import "gorm.io/datatypes"

type CommonWriter struct {
	CommonUUID
	CommonSyncUUID
	WriterThingClass string         `json:"writer_thing_class,omitempty"`
	WriterThingType  string         `json:"writer_thing_type,omitempty"`
	WriterThingUUID  string         `json:"writer_thing_uuid,omitempty"`
	WriterThingName  string         `json:"writer_thing_name,omitempty"`
	DataStore        datatypes.JSON `json:"data_store,omitempty"`
	CommonCreated
}

type Writer struct {
	CommonWriter
	ConsumerUUID string   `json:"consumer_uuid,omitempty" gorm:"TYPE:string REFERENCES consumers;not null;default:null"`
	PresentValue *float64 `json:"present_value,omitempty"`
}

type WriterClone struct { //TODO the WriterClone needs to publish a COV event as for example we have 2x mqtt broker then the cov for a point maybe different when not going over the internet
	CommonWriter
	ProducerUUID      string `json:"producer_uuid" gorm:"TYPE:string REFERENCES producers;not null;default:null"`
	FlowFrameworkUUID string `json:"flow_framework_uuid,omitempty"`
	CommonSourceUUID
}

type SyncWriter struct {
	ProducerUUID      string
	WriterUUID        string
	FlowFrameworkUUID string
}

type WriterBody struct {
	Action   string        `json:"action,omitempty"`
	Priority *Priority     `json:"priority,omitempty"`
	Schedule *ScheduleData `json:"schedule,omitempty"`
}

type WriterBulkBody struct {
	WriterUUID string `json:"writer_uuid,omitempty"`
	WriterBody
}

type WriterActionOutput struct {
	UUID         string          `json:"uuid"`
	Action       string          `json:"action"`
	IsError      bool            `json:"is_error"`
	Message      *string         `json:"message,omitempty"`
	PresentValue *float64        `json:"present_value,omitempty"`
	DataStore    *datatypes.JSON `json:"data_store,omitempty"`
}

type SyncCOV struct {
	Priority *Priority
	Schedule *ScheduleData
}

type SyncWriterAction struct {
	Priority *Priority
	Schedule *ScheduleData
}
