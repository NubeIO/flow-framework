package versions

import "github.com/NubeIO/flow-framework/model"

func GetInitInterfaces() []interface{} {
	return []interface{}{
		&model.LocalStorageFlowNetwork{},
		&model.Alert{},
		&model.Message{},
		&model.PluginConf{},
		&model.Network{},
		&model.Device{},
		&model.Point{},
		&model.FlowNetwork{},
		&model.FlowNetworkClone{},
		&model.Priority{},
		&model.ProducerHistory{},
		&model.ConsumerHistory{},
		&model.Job{},
		&model.Stream{},
		&model.StreamClone{},
		&model.CommandGroup{},
		&model.Producer{},
		&model.Consumer{},
		&model.Writer{},
		&model.WriterClone{},
		&model.Integration{},
		&model.MqttConnection{},
		&model.Schedule{},
		&model.Tag{},
		&model.Block{},
		&model.Connection{},
		&model.BlockStaticRoute{},
		&model.BlockRouteValueNumber{},
		&model.BlockRouteValueString{},
		&model.BlockRouteValueBool{},
		&model.SourceParameter{},
		&model.Link{},
		&model.History{},
		&model.HistoryLog{},
	}
}
