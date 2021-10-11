package model

type CommonFlowNetwork struct {
	CommonUUID
	CommonSyncUUID
	CommonName
	CommonDescription
	GlobalUUID              string    `json:"global_uuid,omitempty" gorm:"unique"`
	ClientId                string    `json:"client_id,omitempty"`
	ClientName              string    `json:"client_name,omitempty"`
	SiteId                  string    `json:"site_id,omitempty"`
	SiteName                string    `json:"site_name,omitempty"`
	DeviceId                string    `json:"device_id,omitempty"`
	DeviceName              string    `json:"device_name,omitempty"`
	IsRemote                *bool     `json:"is_remote,omitempty"`
	FetchHistories          *bool     `json:"fetch_histories,omitempty"`
	FetchHistoriesFrequency int       `json:"fetch_hist_frequency,omitempty"`      //time example 15min
	DeleteHistoriesOnFetch  *bool     `json:"delete_histories_on_fetch,omitempty"` //drop the histories on the producer device on migration
	IsMQTT                  *bool     `json:"is_mqtt,omitempty"`
	FlowHTTPS               *bool     `json:"flow_https,omitempty"`
	FlowIP                  string    `json:"flow_ip,omitempty"`
	FlowPort                int       `json:"flow_port,omitempty"`
	FlowUsername            string    `json:"flow_username,omitempty"`
	FlowPassword            string    `json:"flow_password,omitempty"`
	FlowToken               string    `json:"flow_token,omitempty"`
	Streams                 []*Stream `json:"streams" gorm:"many2many:flow_networks_streams;constraint:OnDelete:CASCADE"`
	CommonCreated
}

type FlowNetwork struct {
	CommonFlowNetwork
	FlowIP   string    `json:"flow_ip,omitempty" gorm:"uniqueIndex:ip_port_composite_index"`
	FlowPort int       `json:"flow_port,omitempty" gorm:"uniqueIndex:ip_port_composite_index"`
	Streams  []*Stream `json:"streams" gorm:"many2many:flow_networks_streams;constraint:OnDelete:CASCADE"`
	Schedules  []*Schedules  `json:"schedules,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

type FlowNetworkClone struct {
	CommonFlowNetwork
	CommonSourceUUID
	FlowIP       string         `json:"flow_ip,omitempty" gorm:"uniqueIndex:ip_port_clone_composite_index"`
	FlowPort     int            `json:"flow_port,omitempty" gorm:"uniqueIndex:ip_port_clone_composite_index"`
	StreamClones []*StreamClone `json:"stream_clones" gorm:"constraint:OnDelete:CASCADE;"`
}
