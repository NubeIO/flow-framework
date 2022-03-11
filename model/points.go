package model

import (
	"github.com/NubeIO/flow-framework/src/poller"
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

type ObjectType string

const (
	//bacnet
	ObjTypeAnalogInput  ObjectType = "analog_input"
	ObjTypeAnalogOutput ObjectType = "analog_output"
	ObjTypeAnalogValue  ObjectType = "analog_value"
	ObjTypeBinaryInput  ObjectType = "binary_input"
	ObjTypeBinaryOutput ObjectType = "binary_output"
	ObjTypeBinaryValue  ObjectType = "binary_value"
	ObjAnalogInput      ObjectType = "analogInput"
	ObjAnalogOutput     ObjectType = "analogOutput"
	ObjAnalogValue      ObjectType = "analogValue"
	ObjBinaryInput      ObjectType = "binaryInput"
	ObjBinaryOutput     ObjectType = "binaryOutput"
	ObjBinaryValue      ObjectType = "binaryValue"

	//modbus
	ObjTypeReadCoil           ObjectType = "read_coil"
	ObjTypeReadCoils          ObjectType = "read_coils"
	ObjTypeReadDiscreteInput  ObjectType = "read_discrete_input"
	ObjTypeReadDiscreteInputs ObjectType = "read_discrete_inputs"
	ObjTypeWriteCoil          ObjectType = "write_coil"
	ObjTypeWriteCoils         ObjectType = "write_coils"
	ObjTypeReadRegister       ObjectType = "read_register"
	ObjTypeReadRegisters      ObjectType = "read_registers"
	ObjTypeReadHolding        ObjectType = "read_holding"
	ObjTypeReadHoldings       ObjectType = "read_holdings"
	ObjTypeWriteHolding       ObjectType = "write_holding"
	ObjTypeWriteHoldings      ObjectType = "write_holdings"
	ObjTypeReadInt16          ObjectType = "read_int_16"
	ObjTypeWriteInt16         ObjectType = "write_int_16"
	ObjTypeReadUint16         ObjectType = "read_uint_16"
	ObjTypeWriteUint16        ObjectType = "write_uint_16"
	ObjTypeReadInt32          ObjectType = "read_int_32"
	ObjTypeReadUint32         ObjectType = "read_uint_32"
	ObjTypeReadFloat32        ObjectType = "read_float_32"
	ObjTypeWriteFloat32       ObjectType = "write_float_32"
	ObjTypeReadFloat64        ObjectType = "read_float_64"
	ObjTypeWriteFloat64       ObjectType = "write_float_64"
)

var ObjectTypesMap = map[ObjectType]int8{
	//bacnet
	ObjTypeAnalogInput: 0, ObjTypeAnalogOutput: 0, ObjTypeAnalogValue: 0,
	ObjTypeBinaryInput: 0, ObjTypeBinaryOutput: 0, ObjTypeBinaryValue: 0,
	ObjAnalogInput: 0, ObjAnalogOutput: 0, ObjAnalogValue: 0,
	ObjBinaryInput: 0, ObjBinaryOutput: 0, ObjBinaryValue: 0,
	//modbus
	ObjTypeReadCoil: 0, ObjTypeReadCoils: 0, ObjTypeReadDiscreteInput: 0,
	ObjTypeReadDiscreteInputs: 0, ObjTypeWriteCoil: 0, ObjTypeWriteCoils: 0,
	ObjTypeReadRegister: 0, ObjTypeReadRegisters: 0, ObjTypeReadHolding: 0,
	ObjTypeReadHoldings: 0, ObjTypeWriteHolding: 0, ObjTypeWriteHoldings: 0,
	ObjTypeReadInt16: 0, ObjTypeWriteInt16: 0,
	ObjTypeReadUint16: 0, ObjTypeWriteUint16: 0,
	ObjTypeReadInt32:   0,
	ObjTypeReadUint32:  0,
	ObjTypeReadFloat32: 0, ObjTypeWriteFloat32: 0,
	ObjTypeReadFloat64: 0, ObjTypeWriteFloat64: 0,
}

type DataType string

const (
	TypeInt16   DataType = "int16"
	TypeUint16  DataType = "uint16"
	TypeInt32   DataType = "int32"
	TypeUint32  DataType = "uint32"
	TypeInt64   DataType = "int64"
	TypeUint64  DataType = "uint64"
	TypeFloat32 DataType = "float32"
	TypeFloat64 DataType = "float64"
)

type ByteOrder string

const (
	ByteOrderLebBew ByteOrder = "leb_bew" //LITTLE_ENDIAN, HIGH_WORD_FIRST
	ByteOrderLebLew ByteOrder = "leb_lew"
	ByteOrderBebLew ByteOrder = "beb_lew"
	ByteOrderBebBew ByteOrder = "beb_bew"
)

type IOType string

const (
	IOTypeRAW           IOType = "raw"
	IOTypeDigital       IOType = "digital"
	IOTypeAToDigital    IOType = "a_to_digital"
	IOTypeVoltageDC     IOType = "voltage_dc"
	IOTypeCurrent       IOType = "current"
	IOTypeThermistor    IOType = "thermistor"
	IOTypeThermistor10K IOType = "thermistor_10k_type_2"
)

type EvalExamples string

const (
	EvalExExample             EvalExamples = "example"
	EvalExBoolInvert          EvalExamples = "bool_invert"
	EvalExCelsiusToFahrenheit EvalExamples = "celsius_to_fahrenheit"
)

type EvalMode string

const (
	EvalModeEnable              EvalMode = "enable"
	EvalModeDisabled            EvalMode = "disabled"
	EvalModeCalcOnOriginalValue EvalMode = "calc_on_original_value"
	EvalModeCalcAfterScale      EvalMode = "calc_after_scale"
)

//Point table
type Point struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonCreated
	CommonThingClass
	CommonThingRef
	CommonThingType
	CommonFault
	PresentValue         *float64            `json:"present_value"` //point value, read only
	OriginalValue        *float64            `json:"original_value"`
	CurrentPriority      *int                `json:"current_priority,omitempty"`
	InSync               *bool               `json:"in_sync"`                    //is set to false when a new value is written from the user example: if its false then modbus would write the new value. if user edits the point it will disable the COV for one time
	WriteValueOnce       *bool               `json:"write_value_once,omitempty"` //when point is used for polling and if it's a writeable point and WriteValueOnce is true then on a successful write it will set the WriteValueOnceSync to true and on the next poll cycle it will not send the write value
	WriteValueOnceSync   *bool               `json:"write_value_once_sync,omitempty"`
	Fallback             *float64            `json:"fallback"`
	DeviceUUID           string              `json:"device_uuid,omitempty" gorm:"TYPE:string REFERENCES devices;not null;default:null"`
	NetworkUUID          string              `json:"network_uuid,omitempty" gorm:"TYPE:string REFERENCES networks;not null;default:null"`
	EnableWriteable      *bool               `json:"writeable,omitempty"`
	IsOutput             *bool               `json:"is_output,omitempty"`
	EvalMode             string              `json:"eval_mode,omitempty"`
	Eval                 string              `json:"eval_expression,omitempty"`
	EvalExample          string              `json:"eval_example,omitempty"`
	COV                  *float64            `json:"cov"`
	ObjectType           string              `json:"object_type,omitempty"`     //binaryInput, coil, if type os input don't return the priority array
	DataType             string              `json:"data_type,omitempty"`       //int16, uint16, float32
	ObjectEncoding       string              `json:"object_encoding,omitempty"` //BEB_LEW bebLew
	IoNumber             string              `json:"io_number,omitempty"`       //DI1,UI1,AO1, temp, pulse, motion
	IoType               string              `json:"io_type,omitempty"`         //0-10dc, 0-40ma, thermistor
	AddressID            *int                `json:"address_id"`                // for example a modbus address or bacnet address
	AddressLength        *int                `json:"address_length"`            // for example a modbus address offset
	AddressUUID          *string             `json:"address_uuid,omitempty"`    // for example a droplet id (so a string)
	NextAvailableAddress *bool               `json:"use_next_available_address,omitempty"`
	Decimal              *uint32             `json:"decimal,omitempty"`
	LimitMin             *float64            `json:"limit_min"`
	LimitMax             *float64            `json:"limit_max"`
	ScaleInMin           *float64            `json:"scale_in_min"`
	ScaleInMax           *float64            `json:"scale_in_max"`
	ScaleOutMin          *float64            `json:"scale_out_min"`
	ScaleOutMax          *float64            `json:"scale_out_max"`
	UnitType             *string             `json:"unit_type,omitempty"` //temperature
	Unit                 *string             `json:"unit,omitempty"`
	UnitTo               *string             `json:"unit_to,omitempty"` //with take the unit and convert to, this would affect the presentValue and the original value will be stored in the raw
	IsProducer           *bool               `json:"is_producer,omitempty"`
	IsConsumer           *bool               `json:"is_consumer,omitempty"`
	ValueRaw             string              `json:"value_raw,omitempty"`
	Priority             *Priority           `json:"priority,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Tags                 []*Tag              `json:"tags,omitempty" gorm:"many2many:points_tags;constraint:OnDelete:CASCADE"`
	WriteMode            poller.WriteMode    `json:"write_mode"`
	WritePollRequired    *bool               `json:"write_required"`
	ReadPollRequired     *bool               `json:"read_required"`
	PollPriority         poller.PollPriority `json:"poll_priority"`
	PollRate             poller.PollRate     `json:"poll_rate"`
	//PollTimer            *time.Timer         `json:"poll_timer,omitempty"`
}

type Priorities struct {
	P1  *float64 `json:"_1"`
	P2  *float64 `json:"_2"`
	P3  *float64 `json:"_3"`
	P4  *float64 `json:"_4"`
	P5  *float64 `json:"_5"`
	P6  *float64 `json:"_6"`
	P7  *float64 `json:"_7"`
	P8  *float64 `json:"_8"`
	P9  *float64 `json:"_9"`
	P10 *float64 `json:"_10"`
	P11 *float64 `json:"_11"`
	P12 *float64 `json:"_12"`
	P13 *float64 `json:"_13"`
	P14 *float64 `json:"_14"`
	P15 *float64 `json:"_15"`
	P16 *float64 `json:"_16"`
}

type Priority struct {
	PointUUID string   `json:"point_uuid,omitempty" gorm:"REFERENCES points;not null;default:null;primaryKey"`
	P1        *float64 `json:"_1"`
	P2        *float64 `json:"_2"`
	P3        *float64 `json:"_3"`
	P4        *float64 `json:"_4"`
	P5        *float64 `json:"_5"`
	P6        *float64 `json:"_6"`
	P7        *float64 `json:"_7"`
	P8        *float64 `json:"_8"`
	P9        *float64 `json:"_9"`
	P10       *float64 `json:"_10"`
	P11       *float64 `json:"_11"`
	P12       *float64 `json:"_12"`
	P13       *float64 `json:"_13"`
	P14       *float64 `json:"_14"`
	P15       *float64 `json:"_15"`
	P16       *float64 `json:"_16"`
}

func (p *Priority) GetHighestPriorityValue() *float64 {
	if p.P1 != nil {
		return p.P1
	}
	if p.P2 != nil {
		return p.P2
	}
	if p.P3 != nil {
		return p.P3
	}
	if p.P4 != nil {
		return p.P4
	}
	if p.P5 != nil {
		return p.P5
	}
	if p.P6 != nil {
		return p.P6
	}
	if p.P7 != nil {
		return p.P7
	}
	if p.P8 != nil {
		return p.P8
	}
	if p.P9 != nil {
		return p.P9
	}
	if p.P10 != nil {
		return p.P10
	}
	if p.P11 != nil {
		return p.P11
	}
	if p.P12 != nil {
		return p.P12
	}
	if p.P13 != nil {
		return p.P13
	}
	if p.P14 != nil {
		return p.P14
	}
	if p.P15 != nil {
		return p.P15
	}
	if p.P16 != nil {
		return p.P16
	}
	return nil
}
