package model

/*

Rubix V2

Network
- Any new network is called a ProducerNetwork for example modbus or lora
- Any ProducerNetwork doesn't need to use point-mapping to get the data into the cloud
- Any ProducerNetwork can be added to a one or more SubscriberNetworks, The ProducerNetwork will keep a leger or which Subscriber are reading/writing to its points

The SubscriberNetwork is the remote rubix device.
So when the ProducerNetwork (the producer network producers data ie: lora) network has a connection with the SubscriberNetwork the ProducerNetwork keeps a ledger of the SubscriberPoints

ProducerNetwork and SubscriberNetwork Jobs
- publish any CRUD updates to all subscribers (ie when a point is deleted or the name is updated)
- publish any COV events

ProducerNetwork settings
- COV will be set in the producer

SubscriberNetwork settings (these settings are not like 2-way meaning that in the SubscriberNetwork if the COV is updated it will not affect the ProducerNetwork setting)
- as this would be considered a normal point in the SubscriberNetwork this point will have all the same settings ie: history, cov and so on

CommandGroup
- is for issuing global schedule writes or global point writes (as in send a value to any point associated with this group)

TimeOverride
- where a point value can be overridden for a duration of time



REST calls
ProducerNetwork
- can call all attributes

*/




type Stream struct {
	CommonUUID
	CommonName
	CommonDescription
	CommonEnable
	CommonCreated
	IsRemote 			bool 	`json:"is_remote"`
	CommonRubixUUID
	Subscriber			[]Subscriber `json:"subscribers" gorm:"constraint:OnDelete:CASCADE;"`
	Subscription		[]Subscriber `json:"subscription" gorm:"constraint:OnDelete:CASCADE;"`
	CommandGroup		string `json:"command_group"`
}