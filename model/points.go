package model

import (
	"gorm.io/datatypes"
)

// TimeOverride TODO add in later
//TimeOverride where a point value can be overridden for a duration of time
type TimeOverride struct {
	PointUUID string `json:"point_uuid" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	StartDate string `json:"start_date"` // START at 25:11:2021:13:00
	EndDate   string `json:"end_date"`   // START at 25:11:2021:13:30
	Value     string `json:"value"`
	Priority  string `json:"priority"`
}

//MathOperation same as in lora and point-server TODO add in later
type MathOperation struct {
	Calc string //x + 1
	X    float64
}

var ObjectTypes = struct {
	//bacnet
	AnalogInput  string
	AnalogOutput string
	AnalogValue  string
	BinaryInput  string
	BinaryOutput string
	BinaryValue  string
	//modbus
	ReadCoil           string
	ReadCoils          string
	ReadDiscreteInput  string
	ReadDiscreteInputs string
	WriteCoil          string
	WriteCoils         string
	ReadRegister       string
	ReadRegisters      string
	ReadInt16          string
	ReadSingleInt16    string
	WriteSingleInt16   string
	ReadUint16         string
	ReadSingleUint16   string
	WriteSingleUint16  string
	ReadInt32          string
	ReadSingleInt32    string
	WriteSingleInt32   string
	ReadUint32         string
	ReadSingleUint32   string
	WriteSingleUint32  string
	ReadFloat32        string
	ReadSingleFloat32  string
	WriteSingleFloat32 string
	ReadFloat64        string
	ReadSingleFloat64  string
	WriteSingleFloat64 string
}{
	//bacnet
	AnalogInput:  "analogInput",
	AnalogOutput: "analogOutput",
	AnalogValue:  "analogValue",
	BinaryInput:  "binaryInput",
	BinaryOutput: "binaryOutput",
	BinaryValue:  "binaryValue",
	//modbus
	ReadCoil:           "readCoil",
	ReadCoils:          "readCoils",
	ReadDiscreteInput:  "readDiscreteInput",
	ReadDiscreteInputs: "readDiscreteInputs",
	WriteCoil:          "writeCoil",
	WriteCoils:         "writeCoils",
	ReadRegister:       "readRegister",
	ReadRegisters:      "readRegisters",
	ReadInt16:          "readInt16",
	ReadSingleInt16:    "readSingleInt16",
	WriteSingleInt16:   "writeSingleInt16",
	ReadUint16:         "readUint16",
	ReadSingleUint16:   "readSingleUint16",
	WriteSingleUint16:  "writeSingleUint16",
	ReadInt32:          "readInt32",
	ReadSingleInt32:    "readSingleInt32",
	WriteSingleInt32:   "writeSingleInt32",
	ReadUint32:         "readUint32",
	ReadSingleUint32:   "readSingleUint32",
	WriteSingleUint32:  "writeSingleUint32",
	ReadFloat32:        "readFloat32",
	ReadSingleFloat32:  "readSingleFloat32",
	WriteSingleFloat32: "writeSingleFloat32",
	ReadFloat64:        "readFloat64",
	ReadSingleFloat64:  "readSingleFloat64",
	WriteSingleFloat64: "writeSingleFloat64",
}

var ObjectEncoding = struct {
	LebBew string //LITTLE_ENDIAN, HIGH_WORD_FIRST
	LebLew string
	BebLew string
	BebBew string
}{
	LebBew: "lebBew",
	LebLew: "lebLew",
	BebLew: "bebLew",
	BebBew: "bebBew",
}

//Point table
type Point struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	PresentValue         float64        `json:"present_value,omitempty"` //point value, read only
	ValueDisplay         string         `json:"value_display,omitempty"` //point value, read only
	ValueOriginal        *float64       `json:"value_original,omitempty"`
	ValueRaw             datatypes.JSON `json:"value_raw,omitempty"`
	CurrentPriority      *int           `json:"current_priority,omitempty"`
	Fallback             float64        `json:"fallback,omitempty"`
	DeviceUUID           string         `json:"device_uuid,omitempty" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	EnableWriteable      *bool          `json:"writeable,omitempty"`
	IsOutput             *bool          `json:"is_output,omitempty"`
	BoolInvert           *bool          `json:"bool_invert,omitempty"`
	COV                  *float32       `json:"cov,omitempty"`
	ObjectType           string         `json:"object_type,omitempty"`    //binaryInput, coil, if type os input dont return the priority array  TODO decide if we just stick to bacnet object types, as a binaryOut is the sample as a coil in modbus
	AddressId            *int           `json:"address_id,omitempty"`     // for example a modbus address or bacnet address
	AddressOffset        *int           `json:"address_offset,omitempty"` // for example a modbus address offset
	AddressUUID          string         `json:"address_uuid,omitempty"`   // for example a droplet id (so a string)
	NextAvailableAddress *bool          `json:"use_next_available_address,omitempty"`
	Decimal              *uint32        `json:"decimal,omitempty"`
	LimitMin             *float64       `json:"limit_min,omitempty"`
	LimitMax             *float64       `json:"limit_max,omitempty"`
	ScaleInMin           *float64       `json:"scale_in_min,omitempty"`
	ScaleInMax           *float64       `json:"scale_in_max,omitempty"`
	ScaleOutMin          *float64       `json:"scale_out_min,omitempty"`
	ScaleOutMax          *float64       `json:"scale_out_max,omitempty"`
	UnitType             string         `json:"unit_type,omitempty"` //temperature
	Unit                 string         `json:"unit,omitempty"`
	UnitTo               string         `json:"unit_to,omitempty"` //with take the unit and convert to, this would affect the presentValue and the original value will be stored in the raw
	CommonThingClass
	CommonThingRef
	CommonThingType
	IsProducer *bool `json:"is_producer"`
	IsConsumer *bool `json:"is_consumer"`
	CommonFault
	Priority *Priority `json:"priority" gorm:"constraint:OnDelete:CASCADE"`
	Tags     []*Tag    `json:"tags" gorm:"many2many:points_tags;constraint:OnDelete:CASCADE"`
}

type Priority struct {
	PointUUID string   `json:"point_uuid,omitempty" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	P1        *float64 `json:"_1,omitempty"` //would be better if we stored the TS and where it was written from, for example from a Remote Producer
	P2        *float64 `json:"_2,omitempty"`
	P3        *float64 `json:"_3,omitempty"`
	P4        *float64 `json:"_4,omitempty"`
	P5        *float64 `json:"_5,omitempty"`
	P6        *float64 `json:"_6,omitempty"`
	P7        *float64 `json:"_7,omitempty"`
	P8        *float64 `json:"_8,omitempty"`
	P9        *float64 `json:"_9,omitempty"`
	P10       *float64 `json:"_10,omitempty"`
	P11       *float64 `json:"_11,omitempty"`
	P12       *float64 `json:"_12,omitempty"`
	P13       *float64 `json:"_13,omitempty"`
	P14       *float64 `json:"_14,omitempty"`
	P15       *float64 `json:"_15,omitempty"`
	P16       *float64 `json:"_16,omitempty"` //removed and added to the point to save one DB write
}
