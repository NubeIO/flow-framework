package lorawan_model

type BasePayload struct {
	ApplicationID   string `json:"applicationID"`
	ApplicationName string `json:"applicationName"`
	DeviceName      string `json:"deviceName"`
	DevEUI          string `json:"devEUI"`
}

//ElsysAPB from elis
type ElsysAPB struct {
	ApplicationID   string `json:"applicationID"`
	ApplicationName string `json:"applicationName"`
	DeviceName      string `json:"deviceName"`
	DevEUI          string `json:"devEUI"`
	RxInfo          []struct {
		GatewayID string  `json:"gatewayID"`
		UplinkID  string  `json:"uplinkID"`
		Name      string  `json:"name"`
		Rssi      int     `json:"rssi"`
		LoRaSNR   float64 `json:"loRaSNR"`
		Location  struct {
			Latitude  int `json:"latitude"`
			Longitude int `json:"longitude"`
			Altitude  int `json:"altitude"`
		} `json:"location"`
	} `json:"rxInfo"`
	TxInfo struct {
		Frequency int `json:"frequency"`
		Dr        int `json:"dr"`
	} `json:"txInfo"`
	Adr    bool   `json:"adr"`
	FCnt   int    `json:"fCnt"`
	FPort  int    `json:"fPort"`
	Data   string `json:"data"`
	Object struct {
		Pressure  int `json:"pressure"`
		PulseAbs  int `json:"pulseAbs"`
		PulseAbs2 int `json:"pulseAbs2"`
	} `json:"object"`
}
