package database

import "github.com/NubeIO/flow-framework/model"

func (d *GormDatabase) GetHistoryInfluxTags(producerUuid string) ([]*model.HistoryInfluxTag, error) {
	var influxHistoryTags []*model.HistoryInfluxTag
	d.DB.Table("flow_network_clones").
		Select("flow_network_clones.client_id, flow_network_clones.client_name, flow_network_clones.site_id, "+
			"flow_network_clones.site_name, flow_network_clones.device_id, flow_network_clones.device_name, "+
			"points.uuid as point_uuid, points.name as point_name, devices.uuid as device_uuid, "+
			"devices.name as device_name, networks.uuid as networks_uuid, networks.name as network_name").
		Joins("INNER JOIN stream_clones ON flow_network_clones.uuid = stream_clones.flow_network_clone_uuid").
		Joins("INNER JOIN consumers ON stream_clones.uuid = consumers.stream_clone_uuid").
		Joins("INNER JOIN writers ON consumers.uuid = writers.consumer_uuid").
		Joins("INNER JOIN points ON points.uuid = writers.writer_thing_uuid").
		Joins("INNER JOIN devices ON devices.uuid = points.device_uuid").
		Joins("INNER JOIN networks ON networks.uuid = devices.network_uuid").
		Where("consumers.producer_uuid = ?", producerUuid).
		Scan(&influxHistoryTags)
	return influxHistoryTags, nil
}
