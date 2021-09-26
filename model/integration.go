package model

//IntegrationCredentials is to be used when a plugin wants to use username and password for an
type IntegrationCredentials struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	Username        string
	Password        string
	Token           string
	HeaderName      string //used for headers if they are needed for auth
	HeaderValue     string
	IntegrationUUID string `json:"integration_uuid" gorm:"TYPE:varchar(255) REFERENCES integrations;null;default:null"`
}

type Integration struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonFault
	IP                     string                   `json:"ip"`
	PORT                   string                   `json:"port"`
	PluginName             string                   `json:"plugin_name"`
	IntegrationType        string                   `json:"integration_type"`
	PluginConfId           string                   `json:"plugin_conf_id" gorm:"TYPE:varchar(255) REFERENCES plugin_confs;not null;default:null"`
	MqttConnection         []MqttConnection         `json:"mqtt_connections" gorm:"constraint:OnDelete:CASCADE;"`
	IntegrationCredentials []IntegrationCredentials `json:"integration_credentials" gorm:"constraint:OnDelete:CASCADE;"`
	CommonCreated
}
