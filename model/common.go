// Package model contains all models used across the flow-framework
package model

import (
	"time"

	"github.com/NubeIO/null"
)

type CommonDescription struct {
	Description string `json:"description,omitempty"`
}

type CommonName struct {
	Name string `json:"name"`
}

type CommonThingClass struct {
	ThingClass string `json:"thing_class,omitempty"`
}

type CommonThingType struct {
	ThingType string `json:"thing_type,omitempty"`
}

type CommonThingObject struct {
	ThingObject string `json:"thing_object,omitempty"`
}

type CommonThingRef struct {
	ThingRef string `json:"thing_reference,omitempty"`
}

type CommonThingUUID struct {
	ThingUUID string `json:"thing_uuid,omitempty"`
}

type CommonNameUnique struct {
	Name string `json:"name"  gorm:"type:varchar(255);unique;not null"`
}

type CommonModulePath struct {
	ModulePath string `json:"module_path"  gorm:"type:varchar(255);unique;not null"`
}

type CommonHelp struct {
	Help string `json:"help"`
}

type CommonIP struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type CommonNodeType struct {
	NodeType string `json:"node_type"`
}

type CommonType struct {
	ObjectType string `json:"object_type"`
}

type CommonAction struct {
	Action string `json:"action"`
}

type CommonEnable struct {
	Enable *bool `json:"enable"`
}

type CommonID struct {
	ID string `json:"id"`
}

type CommonIDUnique struct {
	Name string `json:"id" gorm:"type:varchar(255);unique;not null"`
}

type CommonUUID struct {
	UUID string `json:"uuid" sql:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
}

type CommonSourceUUID struct {
	SourceUUID string `json:"source_uuid" gorm:"uniqueIndex;not null"`
}

type CommonSyncUUID struct {
	SyncUUID string `json:"sync_uuid"`
}

type CommonCreated struct {
	CreatedAt time.Time `json:"created_on,omitempty"`
	UpdatedAt time.Time `json:"updated_on,omitempty"`
}

type CommonHistory struct {
	EnableHistory *bool `json:"history_enable,omitempty"`
}

type CommonValue struct {
	Value    null.Float `json:"value,omitempty"`
	ValueRaw string     `json:"value_raw,omitempty"`
}

type CommonFault struct {
	InFault      bool      `json:"fault,omitempty"`
	MessageLevel string    `json:"message_level,omitempty"`
	MessageCode  string    `json:"message_code,omitempty"`
	Message      string    `json:"message,omitempty"`
	LastOk       time.Time `json:"last_ok,omitempty"`
	LastFail     time.Time `json:"last_fail,omitempty"`
}

type CommonStore struct {
	CommonValue
	CommonFault
}

type CommonProducerPermissions struct {
	Blacklist bool `json:"blacklist"`
	ReadOnly  bool `json:"read_only"`
	AllowCRUD bool `json:"allow_crud"` //not sure if this will be used, but it will allow the producer to update the producer
}

type CommonCurrentWriterUUID struct {
	CurrentWriterUUID string `json:"current_writer_uuid"`
}
