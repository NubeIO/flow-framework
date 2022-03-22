package model

type Integration struct {
	CommonUUID
	CommonNameUnique
	CommonDescription
	CommonEnable
	CommonFault
	IP              string           `json:"ip"`
	PORT            int              `json:"port"`
	Username        string           `json:"username"`
	Password        string           `json:"password"`
	Token           string           `json:"token"`
	PluginName      string           `json:"plugin_name"`
	IntegrationType string           `json:"integration_type"`
	PluginConfId    string           `json:"plugin_conf_id" gorm:"TYPE:varchar(255) REFERENCES plugin_confs;not null;default:null"`
	MqttConnection  []MqttConnection `json:"mqtt_connections" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}
