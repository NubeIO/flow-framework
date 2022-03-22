package lorawan_model

import (
	"github.com/NubeIO/null"
	"time"
)

type Organizations struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id              string    `json:"id"`
		Name            string    `json:"name"`
		DisplayName     string    `json:"displayName"`
		CanHaveGateways bool      `json:"canHaveGateways"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
	} `json:"result"`
}

//Devices GET
type Devices struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		DevEUI                              string    `json:"devEUI"`
		Name                                string    `json:"name"`
		ApplicationID                       string    `json:"applicationID"`
		Description                         string    `json:"description"`
		DeviceProfileID                     string    `json:"deviceProfileID"`
		DeviceProfileName                   string    `json:"deviceProfileName"`
		DeviceStatusBattery                 int       `json:"deviceStatusBattery"`
		DeviceStatusMargin                  int       `json:"deviceStatusMargin"`
		DeviceStatusExternalPowerSource     bool      `json:"deviceStatusExternalPowerSource"`
		DeviceStatusBatteryLevelUnavailable bool      `json:"deviceStatusBatteryLevelUnavailable"`
		DeviceStatusBatteryLevel            int       `json:"deviceStatusBatteryLevel"`
		LastSeenAt                          time.Time `json:"lastSeenAt"`
	} `json:"result"`
}

//Device POST and PUT (put is edit)
type Device struct {
	Device struct {
		DevEUI          string `json:"devEUI"`
		Name            string `json:"name"`
		ApplicationID   string `json:"applicationID"`
		Description     string `json:"description"`
		DeviceProfileID string `json:"deviceProfileID"`
	} `json:"device"`
}

type GetDevice struct {
	Device struct {
		DevEUI            string `json:"devEUI"`
		Name              string `json:"name"`
		ApplicationID     string `json:"applicationID"`
		Description       string `json:"description"`
		DeviceProfileID   string `json:"deviceProfileID"`
		SkipFCntCheck     bool   `json:"skipFCntCheck"`
		ReferenceAltitude int    `json:"referenceAltitude"`
		Variables         struct {
		} `json:"variables"`
		Tags struct {
		} `json:"tags"`
		IsDisabled bool `json:"isDisabled"`
	} `json:"device"`
	LastSeenAt          interface{} `json:"lastSeenAt"`
	DeviceStatusBattery int         `json:"deviceStatusBattery"`
	DeviceStatusMargin  int         `json:"deviceStatusMargin"`
	Location            interface{} `json:"location"`
}

//DeviceActivation used to GET Activation
type DeviceActivation struct {
	DA struct {
		DevAddr     string `json:"devAddr"`
		NwkSEncKey  string `json:"nwkSEncKey"`
		AppSKey     string `json:"appSKey"`
		DevEUI      string `json:"devEUI"`
		FNwkSIntKey string `json:"fNwkSIntKey"`
		SNwkSIntKey string `json:"sNwkSIntKey"`
	} `json:"deviceActivation"`
}

//DeviceActive used to POST active
type DeviceActive struct {
	DActive struct {
		DevEUI      string `json:"devEUI"`
		DevAddr     string `json:"devAddr"`
		AppSKey     string `json:"appSKey"`
		NwkSEncKey  string `json:"nwkSEncKey"`
		SNwkSIntKey string `json:"sNwkSIntKey"`
		FNwkSIntKey string `json:"fNwkSIntKey"`
		FCntUp      int    `json:"fCntUp"`
		NFCntDown   int    `json:"nFCntDown"`
		AFCntDown   int    `json:"aFCntDown"`
	} `json:"deviceActivation"`
}

type DeviceProfiles struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id                string    `json:"id"`
		Name              string    `json:"name"`
		OrganizationID    string    `json:"organizationID"`
		NetworkServerID   string    `json:"networkServerID"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
		NetworkServerName string    `json:"networkServerName"`
	} `json:"result"`
}

type Applications struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id                 string `json:"id"`
		Name               string `json:"name"`
		Description        string `json:"description"`
		OrganizationID     string `json:"organizationID"`
		ServiceProfileID   string `json:"serviceProfileID"`
		ServiceProfileName string `json:"serviceProfileName"`
	} `json:"result"`
}

type ServiceProfiles struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id                string    `json:"id"`
		Name              string    `json:"name"`
		OrganizationID    string    `json:"organizationID"`
		NetworkServerID   string    `json:"networkServerID"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
		NetworkServerName string    `json:"networkServerName"`
	} `json:"result"`
}

type GatewayProfiles struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id                string    `json:"id"`
		Name              string    `json:"name"`
		NetworkServerID   string    `json:"networkServerID"`
		NetworkServerName string    `json:"networkServerName"`
		CreatedAt         time.Time `json:"createdAt"`
		UpdatedAt         time.Time `json:"updatedAt"`
	} `json:"result"`
}

type Users struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id         string      `json:"id"`
		Email      string      `json:"email"`
		SessionTTL int         `json:"sessionTTL"`
		IsAdmin    interface{} `json:"isAdmin"`
		IsActive   interface{} `json:"isActive"`
		CreatedAt  time.Time   `json:"createdAt"`
		UpdatedAt  time.Time   `json:"updatedAt"`
	} `json:"result"`
}

type Gateways struct {
	TotalCount string `json:"totalCount"`
	Result     []struct {
		Id              string    `json:"id"`
		Name            string    `json:"name"`
		Description     string    `json:"description"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
		FirstSeenAt     time.Time `json:"firstSeenAt"`
		LastSeenAt      time.Time `json:"lastSeenAt"`
		OrganizationID  string    `json:"organizationID"`
		NetworkServerID string    `json:"networkServerID"`
		Location        struct {
			Latitude  int    `json:"latitude"`
			Longitude int    `json:"longitude"`
			Altitude  int    `json:"altitude"`
			Source    string `json:"source"`
			Accuracy  int    `json:"accuracy"`
		} `json:"location"`
	} `json:"result"`
}

type Server struct {
	Ip                string `json:"ip"`
	Port              int    `json:"port"`
	DeviceId          string `json:"device_id"`
	LocalObjName      string `json:"local_obj_name"`
	ModelName         string `json:"model_name"`
	VendorId          string `json:"vendor_id"`
	VendorName        string `json:"vendor_name"`
	EnableIpByNicName bool   `json:"enable_ip_by_nic_name"`
	IpByNicName       string `json:"ip_by_nic_name"` //eth0
}

// MqttPayload payload from the bacnet server
type MqttPayload struct {
	Value    null.Float
	Priority int
}
