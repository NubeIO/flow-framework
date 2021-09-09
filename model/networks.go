package model

//SerialNetwork type serial
type SerialNetwork struct {
	Port     string `json:"port"`      //dev/tty/USB0
	BaudRate int    `json:"baud_rate"` //9600
	StopBits int    `json:"stop_bits"`
	Parity   int    `json:"parity"`
	DataBits int    `json:"data_bits"`
	Timeout  int    `json:"timeout"`
}

type IPType struct {
	REST string `json:"rest"`
	UDP  string `json:"udp"`
	MQTT string `json:"mqttClient"`
}

//IPNetwork type ip based network
type IPNetwork struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Token    string `json:"token"`
	IPType
}

type Network struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	CommonCreated
	Manufacture      string           `json:"manufacture"`
	Model            string           `json:"model"`
	NetworkType      string           `json:"network_type"  gorm:"type:varchar(255);not null"`   //lora
	TransportType    string           `json:"transport_type"  gorm:"type:varchar(255);not null"` //serial
	PluginConfId     string           `json:"plugin_conf_id" gorm:"TYPE:varchar(255) REFERENCES plugin_confs;not null;default:null"`
	Device           []Device         `json:"devices" gorm:"constraint:OnDelete:CASCADE"`
	SerialConnection SerialConnection `json:"serial_connection" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IpConnection     IpConnection     `json:"ip_connection" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
