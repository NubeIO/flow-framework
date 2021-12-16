package model

import "time"

type CommonDevice struct {
	Manufacture       string        `json:"manufacture,omitempty"` // nube
	Model             string        `json:"model,omitempty"`       // thml
	AddressId         int           `json:"address_id,omitempty"`  // for example a modbus address or bacnet address
	ZeroMode          *bool         `json:"zero_mode,omitempty"`
	PollDelayPointsMS time.Duration `json:"poll_delay_points_ms"`
	AddressUUID       string        `json:"address_uuid" gorm:"type:varchar(255)"` // AAB1213
	CommonIP
}

type Device struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	CommonThingClass
	CommonThingRef
	CommonThingType
	CommonDevice
	NetworkUUID string   `json:"network_uuid,omitempty" gorm:"TYPE:varchar(255) REFERENCES networks;not null;default:null"`
	Points      []*Point `json:"points,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Tags        []*Tag   `json:"tags,omitempty" gorm:"many2many:devices_tags;constraint:OnDelete:CASCADE"`
}
