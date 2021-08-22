package model



type SubscriptionsType struct {
	Network   		string `json:"network"`
	Job   			string `json:"job"`
	Point   		string `json:"point"`
	Alarm   		string `json:"alarm"`

}

type SubscriptionsApplication struct {
	Local   		 string `json:"local"`
	Remote   		string `json:"remote"`
	Plugin   		string `json:"plugin"`


}


//Subscriptions could be a local network, job or alarm and so on
type Subscriptions struct {
	CommonSubscriptions
	SubscriberType 			string  `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	GatewayUUID     		string `json:"gateway_uuid" gorm:"TYPE:string REFERENCES gateways;not null;default:null"`
	PointUUID    			string  `json:"point_uuid" binding:"required" gorm:"TYPE:varchar(255) REFERENCES points;not null;default:null"`

}