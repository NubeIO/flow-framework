package model

var CommonNaming = struct {
	Plugin            string
	Network           string
	Device            string
	Point             string
	Read              string
	Write             string
	Stream            string
	StreamList        string
	Job               string
	Producer          string
	Consumer          string
	Writer            string
	WriterClone       string
	Alert             string
	Mapping           string
	CommandGroup      string
	Rubix             string
	RubixGlobal       string
	FlowNetwork       string
	RemoteFlowNetwork string
	History           string
	ProducerHistory   string
	Histories         string
	Node              string
	Serial            string
	IP                string
	TransportType     string
}{
	Plugin:            "plugin",
	Network:           "network",
	Device:            "device",
	Point:             "point",
	Read:              "read",
	Write:             "write",
	Stream:            "stream",
	StreamList:        "stream_list",
	Job:               "job",
	Producer:          "producer",
	Consumer:          "consumer",
	Writer:            "writer",
	WriterClone:       "writer_clone",
	Alert:             "alert",
	Mapping:           "mapping",
	CommandGroup:      "command_group",
	Rubix:             "rubix",
	RubixGlobal:       "rubix_global",
	FlowNetwork:       "flow_network",
	RemoteFlowNetwork: "remote_flow_network",
	History:           "history",
	ProducerHistory:   "producer_history",
	Histories:         "histories",
	Node:              "node",
	Serial:            "serial",
	TransportType:     "transport_type",
}

var CommonNamingCommandGroup = struct {
	PointWrite     string
	MasterSchedule string
	SilenceAlarm   string
}{
	PointWrite:     "point_write",
	MasterSchedule: "master_schedule",
	SilenceAlarm:   "silence_alarm",
}

var CommonFaultCode = struct {
	ConfigError      string
	SystemError      string
	PluginNotEnabled string
	Offline          string
	Ok               string
}{
	ConfigError:      "configError",
	SystemError:      "systemError",
	PluginNotEnabled: "pluginNotEnabled",
	Offline:          "offline",
	Ok:               "ok",
}

var MessageLevel = struct {
	Info         string
	Critical     string
	NoneCritical string
	Warning      string
	Fail         string
	Normal       string
}{
	Info:         "info",
	Critical:     "critical",
	NoneCritical: "noneCritical",
	Warning:      "warning",
	Fail:         "fail",
	Normal:       "normal",
}

var CommonFaultMessage = struct {
	ConfigError      string
	SystemError      string
	PluginNotEnabled string
	Offline          string
	NetworkMessage   string
}{
	ConfigError:      "config error",
	SystemError:      "system error",
	PluginNotEnabled: "plugin not enabled or no valid message from the network",
	Offline:          "offline",
	NetworkMessage:   "msg for network valid",
}

var TransType = struct {
	Serial string
	IP     string
}{
	Serial: "serial",
	IP:     "IP",
}

var TransClient = struct {
	Client          string
	Server          string
	WirelessGateway string
	Stream          string
}{
	Client:          "client",
	Server:          "server",
	WirelessGateway: "wireless",
	Stream:          "gateway",
}

var TransProtocol = struct {
	REST         string
	BACnet       string
	Modbus       string
	ModbusMaster string
	MQTT         string
	LoRa         string
	LoRaWAN      string
}{
	REST:         "rest",
	BACnet:       "BACnet",
	Modbus:       "Modbus",
	ModbusMaster: "ModbusMaster",
	MQTT:         "MQTT",
	LoRa:         "LoRa",
	LoRaWAN:      "LoRaWAN",
}

var PointTags = struct {
	RSSI     string
	Voltage  string
	Temp     string
	Humidity string
	Light    string
	Motion   string
}{
	RSSI:     "rssi",
	Voltage:  "voltage",
	Temp:     "temperature",
	Humidity: "humidity",
	Light:    "light",
	Motion:   "motion",
}
