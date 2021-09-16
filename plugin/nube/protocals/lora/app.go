package main

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
	log "github.com/sirupsen/logrus"
	"time"
)

/*
user adds a network
user adds a device
- create device and send plugin the uuid
- ask the plugin do you want to add pre-made points for example
- add points
*/

func SerialOpenAndRead() error {
	s := new(SerialSetting)
	s.SerialPort = "/dev/ttyACM0"
	s.BaudRate = 38400
	sc := New(s)
	err := sc.NewSerialConnection()
	if err != nil {
		return err
	}
	sc.Loop()
	return nil
}

// SerialOpen open serial port
func (i *Instance) SerialOpen() error {
	go func() error {
		err := SerialOpenAndRead()
		if err != nil {
			return err
		}
		return nil
	}()
	log.Info("LORA: open serial port")
	return nil
}

// SerialClose close serial port
func (i *Instance) SerialClose() error {
	err := Disconnect()
	if err != nil {
		return err
	}
	return nil
}

/*
NETWORK
*/
// addPoints close serial port
func (i *Instance) validateNetwork(body *model.Network) error {
	//rules
	// max one network for lora-raw, if user adds a 2nd network it will be put into fault
	// serial port must be set

	return nil
}

/*
POINTS
*/

// addPoints close serial port
func (i *Instance) addPoints(deviceBody *model.Device) (*model.Point, error) {
	p := new(model.Point)
	p.DeviceUUID = deviceBody.UUID
	p.AddressUUID = deviceBody.AddressUUID
	if deviceBody.Model == string(decoder.THLM) {
		for _, e := range THLM {
			p.ThingType = e //temp
			err := i.addPoint(p)
			if err != nil {
				log.Error("LORA: issue on add points", " ", err)
				return nil, err
			}
		}
	}
	return nil, nil

}

// addPoints close serial port
func (i *Instance) addPoint(body *model.Point) error {
	_, err := i.db.CreatePoint(body)
	if err != nil {
		return err
	}
	return nil
}

// updatePoints update the point values
func (i *Instance) updatePoints(deviceBody *model.Device) (*model.Point, error) {
	p := new(model.Point)
	p.UUID = deviceBody.UUID
	code := deviceBody.AddressUUID
	if code == string(decoder.THLM) {
		for _, e := range THLM {
			p.ThingType = e
			err := i.updatePoint(p)
			if err != nil {
				log.Error("LORA: issue on add points", " ", err)
				return nil, err
			}
		}
	}
	return nil, nil

}

// updatePoint by its lora id and type as in temp or lux
func (i *Instance) updatePoint(body *model.Point) error {
	addr := body.AddressUUID
	_, err := i.db.UpdatePointByFieldAndType("address_uuid", addr, body)
	if err != nil {
		return err
	}
	return nil
}

// updatePoint by its lora id
func (i *Instance) devTHLM(pnt *model.Point, value float64) error {
	pnt.PresentValue = value
	pnt.CommonFault.InFault = false
	pnt.CommonFault.MessageLevel = model.MessageLevel.Info
	pnt.CommonFault.MessageCode = model.CommonFaultCode.Ok
	pnt.CommonFault.Message = model.CommonFaultMessage.NetworkMessage
	pnt.CommonFault.LastOk = time.Now().UTC()
	err := i.updatePoint(pnt)
	if err != nil {
		log.Error("LORA issue on update points", " ", err)
	}
	return nil
}

var THLM = []string{"rssi", "voltage", "temperature", "humidity", "light", "motion"}

// PublishSensor close serial port
func (i *Instance) publishSensor(commonSensorData decoder.CommonValues, sensorStruct interface{}) {
	pnt := new(model.Point)
	pnt.AddressUUID = commonSensorData.Id

	if commonSensorData.Sensor == string(decoder.THLM) {
		s := sensorStruct.(decoder.TDropletTHLM)
		for _, e := range THLM {
			switch e {
			case model.PointTags.RSSI:
				f := float64(s.Rssi)
				pnt.ThingType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Voltage:
				f := float64(s.Voltage)
				pnt.ThingType = e //set point type
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}
			case model.PointTags.Temp:
				pnt.ThingType = e //set point type
				err := i.devTHLM(pnt, s.Temperature)
				if err != nil {
					return
				}
			case model.PointTags.Humidity:
				pnt.ThingType = e //set point type
				f := float64(s.Humidity)
				err := i.devTHLM(pnt, f)
				if err != nil {
					return
				}

			}
		}
	}
}
