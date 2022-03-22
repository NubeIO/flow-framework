package decoder

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
)

type TZHTBase struct {
	CommonValues
	PayloadType     string `json:"payload_type"`
	ProtocolVersion uint8  `json:"protocol_version"`
}

type TZHTStaticMin struct {
	LoRaFirmwareMajor       uint8  `json:"lora_firmware_major"`
	LoRaFirmwareMinor       uint8  `json:"lora_firmware_minor"`
	LoRaBuildMajor          uint8  `json:"lora_build_major"`
	LoRaBuildMinor          uint8  `json:"lora_build_minor"`
	SerialNumber            string `json:"serial_number"`
	ModelNumber             string `json:"model_number"`
	ProductNumber           string `json:"product_number"`
	FirmwareVersion         string `json:"firmware_version"`
	CalibrationDate         string `json:"calibration_date"`
	First50LitresDate       string `json:"first_50_litres_date"`
	FilterLogDateInternal   string `json:"filter_log_date_internal"`
	FilterLogLitresInternal int    `json:"filter_log_litres_internal"`
	FilterLogDateExternal   string `json:"filter_log_date_external"`
	FilterLogLitresExternal int    `json:"filter_log_litres_external"`
}

type TZHTStatic struct {
	TZHTBase
	TZHTStaticMin
}

const ZipHTTimerLength = 7

type TZHTTimer struct {
	TimeStart   int  `json:"time_start"`
	TimeStop    int  `json:"time_stop"`
	EnableStart bool `json:"enable_start"`
	EnableStop  bool `json:"enable_stop"`
}

type TZHTWriteMin struct {
	Time                         string                      `json:"time"`
	DispenseTimeBoiling          int                         `json:"dispense_time_boiling"`
	DispenseTimeChilled          int                         `json:"dispense_time_chilled"`
	DispenseTimeSparkling        int                         `json:"dispense_time_sparkling"`
	TemperatureSPBoiling         float32                     `json:"temperature_sp_boiling"`
	TemperatureSPChilled         float32                     `json:"temperature_sp_chilled"`
	TemperatureSPSparkling       float32                     `json:"temperature_sp_sparkling"`
	SleepModeSetting             int                         `json:"sleep_mode_setting"`
	FilterInfoLifeLitresInternal int                         `json:"filter_info_life_litres_internal"`
	FilterInfoLifeMonthsInternal int                         `json:"filter_info_life_months_internal"`
	FilterInfoLifeLitresExternal int                         `json:"filter_info_life_litres_external"`
	FilterInfoLifeMonthsExternal int                         `json:"filter_info_life_months_external"`
	SafetyAllowTapChanges        bool                        `json:"safety_allow_tap_changes"`
	SafetyLock                   bool                        `json:"safety_lock"`
	SafetyHotIsolation           bool                        `json:"safety_hot_isolation"`
	SecurityEnable               bool                        `json:"security_enable"`
	SecurityPin                  string                      `json:"security_pin"`
	Timers                       [ZipHTTimerLength]TZHTTimer `json:"timers"`
}

type TZHTWrite struct {
	TZHTBase
	TZHTWriteMin
}

type TZHTPollMin struct {
	Rebooted                          bool    `json:"rebooted"`
	StaticCOVFlag                     bool    `json:"static_cov_flag"`
	WriteCOVFlag                      bool    `json:"write_cov_flag"`
	SleepModeStatus                   int8    `json:"sleep_mode_status"`
	TemperatureNTCBoiling             float32 `json:"temperature_ntc_boiling"`
	TemperatureNTCChilled             float32 `json:"temperature_ntc_chilled"`
	TemperatureNTCStream              float32 `json:"temperature_ntc_stream"`
	TemperatureNTCCondensor           float32 `json:"temperature_ntc_condensor"`
	UsageEnergyKWh                    float32 `json:"usage_energy_kwh"`
	UsageWaterDeltaDispensesBoiling   int     `json:"usage_water_delta_dispenses_boiling"`
	UsageWaterDeltaDispensesChilled   int     `json:"usage_water_delta_dispenses_chilled"`
	UsageWaterDeltaDispensesSparkling int     `json:"usage_water_delta_dispenses_sparkling"`
	UsageWaterDeltaLitresBoiling      float32 `json:"usage_water_delta_litres_boiling"`
	UsageWaterDeltaLitresChilled      float32 `json:"usage_water_delta_litres_chilled"`
	UsageWaterDeltaLitresSparkling    float32 `json:"usage_water_delta_litres_sparkling"`
	Fault1                            uint8   `json:"fault_1"`
	Fault2                            uint8   `json:"fault_2"`
	Fault3                            uint8   `json:"fault_3"`
	Fault4                            uint8   `json:"fault_4"`
	FilterWarningInternal             bool    `json:"filter_warning_internal"`
	FilterWarningExternal             bool    `json:"filter_warning_external"`
	FilterInfoUsageLitresInternal     int     `json:"filter_info_usage_litres_internal"`
	FilterInfoUsageDaysInternal       int     `json:"filter_info_usage_days_internal"`
	FilterInfoUsageLitresExternal     int     `json:"filter_info_usage_litres_external"`
	FilterInfoUsageDaysExternal       int     `json:"filter_info_usage_days_external"`
}

type TZHTPoll struct {
	TZHTBase
	TZHTPollMin
}

var ZHTFullPointsStruct struct {
	TZHTBase
	TZHTStaticMin
	TZHTWriteMin
	TZHTPollMin
}

func GetPointsStructZHT() interface{} {
	return ZHTFullPointsStruct
}

type TZHTPayloadType int

const (
	ErrorData = iota
	StaticData
	WriteData
	PollData
)

const ZHT_HEX_STR_DATA_START = 14

func DecodeZHT(data string, _ *LoRaDeviceDescription) (*CommonValues, interface{}) {
	bytes := getPayloadBytes(data)
	protocolVersion := getProtocolVersion(data)

	switch pl := getPayloadType(data); pl {
	case StaticData:
		payloadFull := staticPayloadDecoder(bytes)
		payloadFull.PayloadType = "static"
		payloadFull.ProtocolVersion = protocolVersion
		return &payloadFull.CommonValues, payloadFull
	case WriteData:
		payloadFull := writePayloadDecoder(bytes)
		payloadFull.PayloadType = "write"
		payloadFull.ProtocolVersion = protocolVersion
		return &payloadFull.CommonValues, payloadFull
	case PollData:
		payloadFull := pollPayloadDecoder(bytes)
		payloadFull.PayloadType = "poll"
		payloadFull.ProtocolVersion = protocolVersion
		return &payloadFull.CommonValues, payloadFull
	}

	return &CommonValues{}, nil
}

func getPayloadType(data string) TZHTPayloadType {
	plID, _ := strconv.ParseInt(data[14:16], 16, 0)
	return TZHTPayloadType(plID)
}

func CheckPayloadLengthZHT(data string) bool {
	payloadLength := len(data) - 10 // removed addr, nonce and MAC
	payloadLength /= 2
	payloadType := getPayloadType(data)
	dataLength, _ := strconv.ParseInt(data[12:14], 16, 0)
	log.Printf("ZHT dataLength: %d\n", dataLength)

	return (payloadType == StaticData && dataLength == 98 && payloadLength > 98) ||
		(payloadType == WriteData && dataLength == 52 && payloadLength > 52) ||
		(payloadType == PollData && dataLength == 41 && payloadLength > 41)
}

func getPayloadBytes(data string) []byte {
	length, _ := strconv.ParseInt(data[12:14], 16, 0)
	bytes, _ := hex.DecodeString(data[16 : 16+((length-1)*2)])
	return bytes
}

func getProtocolVersion(data string) uint8 {
	v, _ := strconv.ParseInt(data[16:18], 16, 0)
	return uint8(v)
}

func bytesToString(bytes []byte) string {
	str := ""
	for _, b := range bytes {
		if b == 0 {
			break
		}
		str += string(b)
	}
	return str
}

func bytesToDate(bytes []byte) string {
	return fmt.Sprintf("%d/%d/%d", bytes[0], bytes[1], bytes[2])
}

func staticPayloadDecoder(data []byte) TZHTStatic {
	index := 1
	fwMa := data[index]
	index += 1
	fwMi := data[index]
	index += 1
	buildMa := data[index]
	index += 1
	buildMi := data[index]
	index += 1
	sn := bytesToString(data[index : index+15])
	index += 15
	mn := bytesToString(data[index : index+20])
	index += 20
	pn := bytesToString(data[index : index+20])
	index += 20
	fw := bytesToString(data[index : index+20])
	index += 20
	calDate := bytesToDate(data[index : index+3])
	index += 3
	f50lDate := bytesToDate(data[index : index+3])
	index += 3
	filtLogDateInt := bytesToDate(data[index : index+3])
	index += 3
	filtLogLitresInt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	filtLogDateExt := bytesToDate(data[index : index+3])
	index += 3
	filtLogLitresExt := int(binary.LittleEndian.Uint16(data[index : index+2]))
	index += 2
	return TZHTStatic{
		TZHTStaticMin: TZHTStaticMin{
			LoRaFirmwareMajor:       fwMa,
			LoRaFirmwareMinor:       fwMi,
			LoRaBuildMajor:          buildMa,
			LoRaBuildMinor:          buildMi,
			SerialNumber:            sn,
			ModelNumber:             mn,
			ProductNumber:           pn,
			FirmwareVersion:         fw,
			CalibrationDate:         calDate,
			First50LitresDate:       f50lDate,
			FilterLogDateInternal:   filtLogDateInt,
			FilterLogLitresInternal: filtLogLitresInt,
			FilterLogDateExternal:   filtLogDateExt,
			FilterLogLitresExternal: filtLogLitresExt,
		},
	}
}

func writePayloadDecoder(data []byte) TZHTWrite {
	index := 1
	time := fmt.Sprintf("%d", binary.LittleEndian.Uint32(data[index:index+4]))
	index += 4
	dispB := int(data[index])
	index += 1
	dispC := int(data[index])
	index += 1
	dispS := int(data[index])
	index += 1
	tempSpB := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempSpC := float32(int(data[index]))
	index += 1
	tempSpS := float32(int(data[index]))
	index += 1
	sm := int(data[index])
	index += 1
	filLyfLtrInt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	filLyfMnthInt := int(data[index])
	index += 1
	filLyfLtrExt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	filLyfMnthExt := int(data[index])
	index += 1
	sfTap := (data[index]>>2)&1 == 1
	sfL := (data[index]>>1)&1 == 1
	sfHi := (data[index]>>0)&1 == 1
	index += 1
	secUI16 := binary.LittleEndian.Uint16(data[index : index+2])
	secEn := secUI16 >= 10000
	secPin := fmt.Sprintf("%.4d", (secUI16 % 10000))
	index += 2

	var timers [ZipHTTimerLength]TZHTTimer
	var u16 uint16
	for i := 0; i < ZipHTTimerLength; i++ {
		u16 = binary.LittleEndian.Uint16(data[index : index+2])
		timers[i].TimeStart = int(u16 % 10000)
		timers[i].EnableStart = u16 >= 10000
		index += 2
		u16 = binary.LittleEndian.Uint16(data[index : index+2])
		timers[i].TimeStop = int(u16 % 10000)
		timers[i].EnableStop = u16 >= 10000
		index += 2
	}

	return TZHTWrite{
		TZHTWriteMin: TZHTWriteMin{
			Time:                         time,
			DispenseTimeBoiling:          dispB,
			DispenseTimeChilled:          dispC,
			DispenseTimeSparkling:        dispS,
			TemperatureSPBoiling:         tempSpB,
			TemperatureSPChilled:         tempSpC,
			TemperatureSPSparkling:       tempSpS,
			SleepModeSetting:             sm,
			FilterInfoLifeLitresInternal: int(filLyfLtrInt),
			FilterInfoLifeMonthsInternal: filLyfMnthInt,
			FilterInfoLifeLitresExternal: int(filLyfLtrExt),
			FilterInfoLifeMonthsExternal: filLyfMnthExt,
			SafetyAllowTapChanges:        sfTap,
			SafetyLock:                   sfL,
			SafetyHotIsolation:           sfHi,
			SecurityEnable:               secEn,
			SecurityPin:                  secPin,
			Timers:                       timers,
		},
	}
}

func pollPayloadDecoder(data []byte) TZHTPoll {
	index := 1
	rebooted := (data[index]>>5)&1 == 1
	sCov := (data[index]>>6)&1 == 1
	wCov := (data[index]>>7)&1 == 1
	sms := int8((data[index]) & 0x3F)
	index += 1
	tempB := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempC := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempS := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	tempCond := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	f1 := data[index]
	index += 1
	f2 := data[index]
	index += 1
	f3 := data[index]
	index += 1
	f4 := data[index]
	index += 1
	kwh := float32(binary.LittleEndian.Uint32(data[index:index+4])) * 0.1
	index += 4
	dltDispB := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltDispC := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltDispS := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	dltLtrB := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	dltLtrC := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	dltLtrS := float32(binary.LittleEndian.Uint16(data[index:index+2])) / 10
	index += 2
	fltrWrnInt := (data[index]>>0)&1 == 1
	fltrWrnExt := (data[index]>>1)&1 == 1
	index += 1
	fltrNfoUseLtrInt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	fltrNfoUseDayInt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	fltrNfoUseLtrExt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2
	fltrNfoUseDayExt := binary.LittleEndian.Uint16(data[index : index+2])
	index += 2

	return TZHTPoll{
		TZHTPollMin: TZHTPollMin{
			Rebooted:                          rebooted,
			StaticCOVFlag:                     sCov,
			WriteCOVFlag:                      wCov,
			SleepModeStatus:                   sms,
			TemperatureNTCBoiling:             tempB,
			TemperatureNTCChilled:             tempC,
			TemperatureNTCStream:              tempS,
			TemperatureNTCCondensor:           tempCond,
			UsageEnergyKWh:                    kwh,
			UsageWaterDeltaDispensesBoiling:   int(dltDispB),
			UsageWaterDeltaDispensesChilled:   int(dltDispC),
			UsageWaterDeltaDispensesSparkling: int(dltDispS),
			UsageWaterDeltaLitresBoiling:      dltLtrB,
			UsageWaterDeltaLitresChilled:      dltLtrC,
			UsageWaterDeltaLitresSparkling:    dltLtrS,
			Fault1:                            f1,
			Fault2:                            f2,
			Fault3:                            f3,
			Fault4:                            f4,
			FilterWarningInternal:             fltrWrnInt,
			FilterWarningExternal:             fltrWrnExt,
			FilterInfoUsageLitresInternal:     int(fltrNfoUseLtrInt),
			FilterInfoUsageDaysInternal:       int(fltrNfoUseDayInt),
			FilterInfoUsageLitresExternal:     int(fltrNfoUseLtrExt),
			FilterInfoUsageDaysExternal:       int(fltrNfoUseDayExt),
		},
	}
}
