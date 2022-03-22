package model

import "time"

type CommonDevice struct {
	Manufacture       string        `json:"manufacture,omitempty"` // nube
	Model             string        `json:"model,omitempty"`       // thml
	AddressId         int           `json:"address_id,omitempty"`  // for example a modbus address or bacnet address
	ZeroMode          *bool         `json:"zero_mode,omitempty"`   //if true means read from address 0 if false will read at 1
	PollDelayPointsMS time.Duration `json:"poll_delay_points_ms"`
	AddressUUID       *string       `json:"address_uuid" gorm:"type:varchar(255);"` // AAB1213
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
	DeviceMac      *int           `json:"device_mac,omitempty"`
	DeviceObjectId *int           `json:"device_object_id,omitempty"`
	NetworkNumber  *int           `json:"network_number,omitempty"` //bacnet network number
	DeviceMask     *int           `json:"device_mask,omitempty"`
	TypeSerial     *bool          `json:"type_serial,omitempty"`
	TransportType  string         `json:"transport_type,omitempty"` //serial, ip
	SupportsRpm    *bool          `json:"supports_rpm,omitempty"`   //bacnet support read property multiple
	SupportsWpm    *bool          `json:"supports_wpm,omitempty"`   //bacnet support write property multiple
	NetworkUUID    string         `json:"network_uuid,omitempty" gorm:"TYPE:varchar(255) REFERENCES networks;not null;default:null"`
	Points         []*Point       `json:"points,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Tags           []*Tag         `json:"tags,omitempty" gorm:"many2many:devices_tags;constraint:OnDelete:CASCADE"`
	FastPollRate   *time.Duration `json:"fast_poll_rate"`
	NormalPollRate *time.Duration `json:"normal_poll_rate"`
	SlowPollRate   *time.Duration `json:"slow_poll_rate"`
}
