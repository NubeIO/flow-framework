package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/utils"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
)

func main() {

	//MAKE CLIENT
	client := resty.New()
	address := "http://0.0.0.0"
	port := "1660"
	url := fmt.Sprintf("%s:%s", address, port)
	urlClient := fmt.Sprintf("%s/%s", url, "client")
	//MAKE TOKEN
	getToken, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"name":"admin"}`).
		SetBasicAuth("admin", "admin").
		Post(urlClient)
	if err != nil {
		log.Println("getToken err:", err, getToken.Status())
	}
	//log.Println("getToken:", getToken, "status", getToken.Status())

	r := gjson.Get(string(getToken.Body()), "token")
	token := r.Str

	//GET TOKEN
	urlUser := fmt.Sprintf("%s/%s", url, "user")
	user, err := client.NewRequest().
		SetHeader("Authorization", token).
		Get(urlUser)
	if err != nil {
		log.Println("user err:", err, user.Status())
	}
	//log.Println("user:", user, "status", user.Status())

	name, _ := utils.MakeUUID()
	name = fmt.Sprintf("name_%s", name)

	//ADD NETWORK
	urlNetwork := fmt.Sprintf("%s/%s", url, "api/networks")
	addNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetBody(map[string]interface{}{"name": name, "description": "description"}).
		Post(urlNetwork)
	if err != nil {
		log.Println("addNetwork err:", err, addNetwork.Status())
	}
	//log.Println("addNetwork:", addNetwork, "status", addNetwork.Status())

	//GET NETWORK
	urlNetworkUUID := fmt.Sprintf("%s/%s", url, "api/networks/{uuid}")
	r = gjson.Get(string(addNetwork.Body()), "uuid")
	getNetworkUUID := r.Str
	log.Println("getNetworkUUID:", getNetworkUUID)
	getNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetPathParams(map[string]string{
			"uuid": getNetworkUUID,
		}).
		Get(urlNetworkUUID)
	if err != nil {
		log.Println("addNetwork err:", err, getNetwork.Status())
	}
	//log.Println("addNetwork:", getNetwork, "status", getNetwork.Status())

	//EDIT NETWORK
	log.Println("getNetworkUUID:", getNetworkUUID, token)
	editNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetBody(map[string]interface{}{"name": "new_name_" + name}).
		SetPathParams(map[string]string{
			"uuid": getNetworkUUID,
		}).
		Patch(urlNetworkUUID)
	if err != nil {
		log.Println("editNetwork err:", err, editNetwork.Status())
	}
	//log.Println("editNetwork:", editNetwork, "status", editNetwork.Status())


	//ADD DEVICE
	urlDevice := fmt.Sprintf("%s/%s", url, "api/devices")
	addDevice, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetBody(map[string]interface{}{"name": name, "description": "description", "network_uuid": getNetworkUUID}).
		Post(urlDevice)
	if err != nil {
		log.Println("addDevice err:", err, addDevice.Status())
	}
	//log.Println("addDevice:", addDevice, "status", addDevice.Status())

	//GET DEVICES
	getDevices, err := client.NewRequest().
		SetHeader("Authorization", token).
		Get(urlDevice)
	if err != nil {
		log.Println("getDevices err:", err, getDevices.Status())
	}
	//log.Println("getDevices:", getDevices, "status", getDevices.Status())


	//ADD POINT
	r = gjson.Get(string(addDevice.Body()), "uuid")
	getDeviceUUID := r.Str
	urlPoints := fmt.Sprintf("%s/%s", url, "api/points")
	addPoint, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetBody(map[string]interface{}{"name": name, "description": "description", "device__uuid": getDeviceUUID}).
		Post(urlDevice)
	if err != nil {
		log.Println("addPoint err:", err, addPoint.Status())
	}
	//log.Println("addDevice:", addDevice, "status", addDevice.Status())


	//GET POINTS
	getPoints, err := client.NewRequest().
		SetHeader("Authorization", token).
		Get(urlPoints)
	if err != nil {
		log.Println("getPoints err:", err, getPoints.Status())
	}

	//GET POINT
	r = gjson.Get(string(addPoint.Body()), "uuid")
	getPointUUID := r.Str
	urlPoint := fmt.Sprintf("%s/%s", url, "api/points/{uuid}")
	getPoint, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetPathParams(map[string]string{
			"point_uuid": getPointUUID,
		}).
		Get(urlPoint)
	if err != nil {
		log.Println("getPoint err:", err, getPoint.Status())
	}


	//urlNetworkUUID := fmt.Sprintf("%s/%s", url, "api/networks/{uuid}")
	//r = gjson.Get(string(addNetwork.Body()), "uuid")
	//getNetworkUUID := r.Str
	//log.Println("getNetworkUUID:", getNetworkUUID)
	//getNetwork, err := client.NewRequest().
	//	SetHeader("Authorization", token).
	//	SetPathParams(map[string]string{
	//		"uuid": getNetworkUUID,
	//	}).
	//	Get(urlNetworkUUID)
	//if err != nil {
	//	log.Println("addNetwork err:", err, getNetwork.Status())
	//}

	if getToken.Status() == "200 OK" {
		fmt.Println("getToken", "PASS")
	} else {
		fmt.Println("getToken", "FAIL")
	}
	if user.Status() == "200 OK" {
		fmt.Println("user", "PASS")
	} else {
		fmt.Println("user", "FAIL")
	}
	if addNetwork.Status() == "200 OK" {
		fmt.Println("addNetwork", "PASS")
	} else {
		fmt.Println("addNetwork", "FAIL", addNetwork.Status())
	}
	if getNetwork.Status() == "200 OK" {
		fmt.Println("getNetwork", "PASS")
	} else {
		fmt.Println("getNetwork", "FAIL", getNetwork.Status())
	}
	if editNetwork.Status() == "200 OK" {
		fmt.Println("editNetwork", "PASS")
	} else {
		fmt.Println("editNetwork", "FAIL", editNetwork.Status())
	}

	if getDevices.Status() == "200 OK" {
		fmt.Println("getDevices", "PASS")
	} else {
		fmt.Println("getDevices", "FAIL", getDevices.Status())
	}

	if addDevice.Status() == "200 OK" {
		fmt.Println("addDevice", "PASS")
	} else {
		fmt.Println("addDevice", "FAIL", addDevice.Status())
	}

	if addPoint.Status() == "200 OK" {
		fmt.Println("addPoint", "PASS")
	} else {
		fmt.Println("addPoint", "FAIL", addPoint.Status())
	}

	if getPoint.Status() == "200 OK" {
		fmt.Println("getPoint", "PASS")
	} else {
		fmt.Println("getPoint", "FAIL", getPoint.Status())
	}

	if getPoints.Status() == "200 OK" {
		fmt.Println("getPoints", "PASS")
	} else {
		fmt.Println("getPoints", "FAIL", getPoints.Status())
	}

	//DELETE NETWORK
	log.Println("getNetworkUUID:", getNetworkUUID)
	deleteNetwork, err := client.NewRequest().
		SetHeader("Authorization", token).
		SetPathParams(map[string]string{
			"uuid": getNetworkUUID,
		}).
		Delete(urlNetworkUUID)
	if err != nil {
		log.Println("deleteNetwork err:", err, deleteNetwork.Status())
	}
	//log.Println("deleteNetwork:", deleteNetwork, "status", deleteNetwork.Status())

	if deleteNetwork.Status() == "200 OK" {
		fmt.Println("deleteNetwork", "PASS")
	} else {
		fmt.Println("deleteNetwork", "FAIL", getPoints.Status())
	}


}
