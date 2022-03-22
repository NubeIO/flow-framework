package model

type HistoryInfluxTag struct {
	ClientId         string `json:"client_id"`
	ClientName       string `json:"client_name"`
	SiteId           string `json:"site_id"`
	SiteName         string `json:"site_name"`
	DeviceId         string `json:"device_id"`
	DeviceName       string `json:"device_name"`
	RubixPointUUID   string `json:"rubix_point_uuid"`
	RubixPointName   string `json:"rubix_point_name"`
	RubixDeviceUUID  string `json:"rubix_device_uuid"`
	RubixDeviceName  string `json:"rubix_device_name"`
	RubixNetworkUUID string `json:"rubix_network_uuid"`
	RubixNetworkName string `json:"rubix_network_name"`
}
