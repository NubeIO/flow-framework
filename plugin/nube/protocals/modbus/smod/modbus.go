package smod

import (
	"github.com/grid-x/modbus"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type RegType uint
type Endianness uint
type WordOrder uint
type Error string

// Error implements the error interface.
func (e Error) Error() (s string) {
	s = string(e)
	return
}

const (
	ParityNone uint = 0
	ParityEven uint = 1
	ParityOdd  uint = 2

	HoldingRegister RegType = 0
	InputRegister   RegType = 1

	// BigEndian endianness of 16-bit registers
	BigEndian    Endianness = 1
	LittleEndian Endianness = 2

	// HighWordFirst word order of 32-bit registers
	HighWordFirst WordOrder = 1
	LowWordFirst  WordOrder = 2

	// errors
	ErrConfigurationError      Error = "configuration error"
	ErrRequestTimedOut         Error = "request timed out"
	ErrIllegalFunction         Error = "illegal function"
	ErrIllegalDataAddress      Error = "illegal data address"
	ErrIllegalDataValue        Error = "illegal data value"
	ErrServerDeviceFailure     Error = "server device failure"
	ErrAcknowledge             Error = "request acknowledged"
	ErrServerDeviceBusy        Error = "server device busy"
	ErrMemoryParityError       Error = "memory parity error"
	ErrGWPathUnavailable       Error = "gateway path unavailable"
	ErrGWTargetFailedToRespond Error = "gateway target device failed to respond"
	ErrBadCRC                  Error = "bad crc"
	ErrShortFrame              Error = "short frame"
	ErrProtocolError           Error = "protocol error"
	ErrBadUnitId               Error = "bad unit id"
	ErrBadTransactionId        Error = "bad transaction id"
	ErrUnknownProtocolId       Error = "unknown protocol identifier"
	ErrUnexpectedParameters    Error = "unexpected parameters"
)

type ModbusClient struct {
	Client           modbus.Client
	RTUClientHandler *modbus.RTUClientHandler
	TCPClientHandler *modbus.TCPClientHandler
	Endianness       Endianness
	WordOrder        WordOrder
	RegType          RegType
	DeviceZeroMode   bool
	Debug            bool
}

func byteArrayToBoolArray(ba []byte) []bool {
	var s []bool
	for _, b := range ba {
		for _, c := range strconv.FormatUint(uint64(b), 2) {
			s = append(s, c == []rune("1")[0])
		}
	}
	return s
}

func convert(data []byte) []bool {
	res := make([]bool, len(data)*8)
	for i := range res {
		res[i] = data[i/8]&(0x80>>byte(i&0x7)) != 0
	}
	return res
}

//SetEncoding Sets the encoding (endianness and word ordering) of subsequent requests.
func (mc *ModbusClient) SetEncoding(endianness Endianness, wordOrder WordOrder) (err error) {

	if endianness != BigEndian && endianness != LittleEndian {
		log.Errorf("unknown endianness value %v", endianness)
		err = ErrUnexpectedParameters
		return
	}

	if wordOrder != HighWordFirst && wordOrder != LowWordFirst {
		log.Errorf("unknown word order value %v", wordOrder)
		err = ErrUnexpectedParameters
		return
	}
	mc.Endianness = endianness
	mc.WordOrder = wordOrder
	return
}

//ReadCoils Reads multiple coils (function code 01).
func (mc *ModbusClient) ReadCoils(addr uint16, quantity uint16) (raw []byte, out float64, err error) {
	raw, err = mc.Client.ReadCoils(addr, quantity)
	if err != nil {
		log.Errorf("modbus-function: failed to ReadCoils: %v\n", err)
		return
	}
	out = float64(raw[0])
	return
}

//ReadDiscreteInputs Reads multiple Discrete Input Registers (function code 02).
func (mc *ModbusClient) ReadDiscreteInputs(addr uint16, quantity uint16) (raw []byte, out float64, err error) {
	raw, err = mc.Client.ReadDiscreteInputs(addr, quantity)
	if err != nil {
		log.Errorf("modbus-function: failed to ReadDiscreteInputs: %v\n", err)
		return
	}
	out = float64(raw[0])
	return
}

//ReadInputRegisters Reads multiple Input Registers (function code 02).
func (mc *ModbusClient) ReadInputRegisters(addr uint16, quantity uint16) (raw []byte, out float64, err error) {
	raw, err = mc.Client.ReadInputRegisters(addr, quantity)
	if err != nil {
		log.Errorf("modbus-function: failed to ReadInputRegisters: %v\n", err)
		return
	}
	// decode payload bytes as uint16s
	decode := bytesToUint16s(mc.Endianness, raw)
	if len(decode) >= 0 {
		out = float64(decode[0])
	}
	return
}

//ReadHoldingRegisters Reads Holding Registers (function code 02).
func (mc *ModbusClient) ReadHoldingRegisters(addr uint16, quantity uint16) (raw []byte, out float64, err error) {
	raw, err = mc.Client.ReadHoldingRegisters(addr, quantity)
	if err != nil {
		log.Errorf("modbus-function: failed to ReadHoldingRegisters: %v\n", err)
		return
	}
	// decode payload bytes as uint16s
	decode := bytesToUint16s(mc.Endianness, raw)
	if len(decode) >= 0 {
		out = float64(decode[0])
	}
	return
}

//ReadFloat32s Reads multiple 32-bit float registers.
func (mc *ModbusClient) ReadFloat32s(addr uint16, quantity uint16, regType RegType) (raw []float32, err error) {
	var mbPayload []byte
	// read 2 * quantity uint16 registers, as bytes
	if regType == HoldingRegister {
		mbPayload, err = mc.Client.ReadHoldingRegisters(addr, quantity*2)
		if err != nil {
			return
		}
	} else {
		mbPayload, err = mc.Client.ReadInputRegisters(addr, quantity*2)
		if err != nil {
			return
		}
	}
	// decode payload bytes as float32s
	raw = bytesToFloat32s(mc.Endianness, mc.WordOrder, mbPayload)
	return
}

//ReadFloat32 Reads a single 32-bit float register.
func (mc *ModbusClient) ReadFloat32(addr uint16, regType RegType) (raw []float32, out float64, err error) {
	raw, err = mc.ReadFloat32s(addr, 1, regType)
	if err != nil {
		log.Errorf("modbus-function: failed to ReadFloat32: %v\n", err)
		return
	}
	out = float64(raw[0])
	return
}

//ReadFloat64s Reads multiple 64-bit float registers.
func (mc *ModbusClient) ReadFloat64s(addr uint16, quantity uint16, regType RegType) (raw []float64, err error) {
	var mbPayload []byte
	// read 2 * quantity uint16 registers, as bytes
	if regType == HoldingRegister {
		mbPayload, err = mc.Client.ReadHoldingRegisters(addr, quantity*2)
		if err != nil {
			return
		}
	} else {
		mbPayload, err = mc.Client.ReadInputRegisters(addr, quantity*2)
		if err != nil {
			return
		}
	}
	// decode payload bytes as float32s
	raw = bytesToFloat64s(mc.Endianness, mc.WordOrder, mbPayload)

	return
}

//ReadFloat64 Reads a single 64-bit float register.
func (mc *ModbusClient) ReadFloat64(addr uint16, regType RegType) (raw []float64, out float64, err error) {
	raw, err = mc.ReadFloat64s(addr, 1, regType)
	if err != nil {
		log.Errorf("modbus-function: failed to ReadFloat32: %v\n", err)
		return
	}
	out = raw[0]
	return
}

//WriteFloat32 Writes a single 32-bit float register.
func (mc *ModbusClient) WriteFloat32(addr uint16, value float64) (raw []byte, out float64, err error) {
	raw, err = mc.Client.WriteMultipleRegisters(addr, 2, float32ToBytes(mc.Endianness, mc.WordOrder, float32(value)))
	if err != nil {
		log.Errorf("modbus-function: failed to WriteFloat32: %v\n", err)
		return
	}
	out = value
	return
}

//WriteSingleRegister write one register
func (mc *ModbusClient) WriteSingleRegister(addr uint16, value uint16) (raw []byte, out float64, err error) {
	raw, err = mc.Client.WriteSingleRegister(addr, value)
	if err != nil {
		log.Errorf("modbus-function: failed to WriteSingleRegister: %v\n", err)
		return
	}
	out = float64(value)
	return
}

//WriteCoil Writes a single coil (function code 05)
func (mc *ModbusClient) WriteCoil(addr uint16, value uint16) (values []byte, out float64, err error) {
	values, err = mc.Client.WriteSingleCoil(addr, value)
	if err != nil {
		log.Errorf("modbus-function: failed to WriteCoil: %v\n", err)
		return
	}
	if value == 0 {
		out = 0
	} else {
		out = 1
	}
	return
}

//func main() {
//	fmt.Println(12132123)
//
//	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
//	handler.BaudRate = 38400
//	handler.DataBits = 8
//	handler.Parity = "N"
//	handler.StopBits = 1
//	handler.SlaveID = 1
//	handler.Timeout = 5 * time.Second
//
//	handler.Connect()
//	defer handler.Close()
//
//	client := modbus.NewClient(handler)
//	var c ModbusClient
//	c.Client = client
//	c.RegType = HoldingRegister
//	c.Endianness = BigEndian
//	c.WordOrder = LowWordFirst
//
//	coil, err := c.ReadCoils(0, 2)
//	fmt.Println(coil)
//	fmt.Println(err)
//
//	f, err := c.ReadFloat32(0, 2)
//	fmt.Println(f)
//	fmt.Println(err)
//
//	if err != nil {
//		return
//	}
//
//}
