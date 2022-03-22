package decoder

import (
	"fmt"
	"testing"
)

//OLDER DROPLET AB payload
//20ABBC903E083B27457200ED000000574D00  20ABBC90
//11AA0203000000000D2E2D95000000E55900  11AA0203
//19ABAA516D084127417400EB000000944400  19ABAA51

//msg.nodeId = hexString.substring(0, 8);
//msg.temp = parseInt("0x" + hexString.substring(10, 12) + hexString.substring(8, 10)) / 100;
//msg.pressure = parseInt("0x" + hexString.substring(14, 16) + hexString.substring(12, 14)) / 10;
//msg.humidity = parseInt("0x" + hexString.substring(16, 18)) % 128;
//msg.movement = parseInt("0x" + hexString.substring(16, 18)) > 127 ? true : false;
//msg.light = parseInt("0x" + hexString.substring(20, 22) + hexString.substring(18, 20));
//msg.voltage = parseInt("0x" + hexString.substring(22, 24)) / 50;
//msg.checksum = "0x" + hexString.substring(30, 32);
//msg.rssi = parseInt("0x" + hexString.substring(32, 34)) * -1;
//msg.snr = parseInt("0x" + hexString.substring(34, 36)) / 10;

//NEW
//E3B2CCD712085C264C0000E728A6C3CCAE0
//74B2BE5925065526280000DF36031DAE3D00
//80B2071E93065926300000EDB1C090FD5B00
//7FAA88D600000009D8000000000000025500 //7FAA88D6 ME
func TestSum(t *testing.T) {

	common, fullData := DecodePayload("20ABBC903E083B27457200ED000000574D00")
	fmt.Println("OLD", common, fullData)

	common, fullData = DecodePayload("11AA0203000000000D2E2D95000000E55900")
	fmt.Println("OLD", common, fullData)

	common, fullData = DecodePayload("19ABAA516D084127417400EB000000944400")
	fmt.Println("OLD", common, fullData)

	common, fullData = DecodePayload("E3B2CCD712085C264C0000E728A6C3CCAE0")
	fmt.Println("NEW", common, fullData)

	common, fullData = DecodePayload("74B2BE5925065526280000DF36031DAE3D00")
	fmt.Println("NEW", common, fullData)

	common, fullData = DecodePayload("80B2071E93065926300000EDB1C090FD5B00")
	fmt.Println("NEW", common, fullData)

	common, fullData = DecodePayload("7FAA88D600000009D8000000000000025500")
	fmt.Println("NEW", common, fullData)

}
