package model

type CommonDevice struct {
	Manufacture string `json:"manufacture"`  // nube
	Model       string `json:"model"`        // thml
	AddressId   int    `json:"address_id"`   // for example a modbus address or bacnet address
	AddressUUID string `json:"address_uuid"` // AAB1213

}

type Device struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	CommonThingClass //point, job
	CommonThingType  // for example temp, rssi, voltage
	CommonDevice
	NetworkUUID string  `json:"network_uuid" gorm:"TYPE:varchar(255) REFERENCES networks;not null;default:null"`
	Point       []Point `json:"points" gorm:"constraint:OnDelete:CASCADE"`
}
