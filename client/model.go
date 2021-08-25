package client

type Token struct {
	ID  	int  	`json:"id"`
	Token  	string 	`json:"token"`
	UUID  	string 	`json:"uuid"`
}


type ResponseBody struct {
	Response ResponseCommon 	`json:"response"`
	Status     string 			`json:"status"`
	Count     string 			`json:"count"`
}

type ResponseCommon struct {

	UUID  			string 	`json:"uuid"`
	Name  			string 	`json:"name"`
	NetworkUUID  	string 	`json:"network_uuid"`
	DeviceUUID  	string 	`json:"device_uuid"`
	PointUUID  		string 	`json:"point_uuid"`
	StreamUUID  	string 	`json:"stream_uuid"`
	GlobalUUID  	string 	`json:"global_uuid"`


}


type Stream struct {
	Name  		string 	`json:"name"`
	IsRemote  	bool 	`json:"is_remote"`
}
type Subscriber struct {
	Name                  	string `json:"name"`
	Enable                	bool   `json:"enable"`
	SubscriberType        	string `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	StreamUUID  	string 	`json:"stream_uuid"`
	FromUUID 				string `json:"from_uuid"`
	ToUUID 					string `json:"to_uuid"`
	IsRemote 				bool 	`json:"is_remote"`
	RemoteRubixUUID			string 	`json:"remote_rubix_uuid"`
}


type Subscription struct {
	Name                  	string `json:"name"`
	Enable                	bool   `json:"enable"`
	SubscriberType        	string `json:"subscriber_type"`
	SubscriberApplication 	string `json:"subscriber_application"`
	StreamUUID  	string 	`json:"stream_uuid"`
	ToUUID 					string `json:"to_uuid"`
	IsRemote 				bool 	`json:"is_remote"`
	RemoteRubixUUID			string 	`json:"remote_rubix_uuid"`
}
