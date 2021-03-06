package database

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	"gorm.io/gorm"
	"strings"
)

func (d *GormDatabase) buildFlowNetworkQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithStreams {
		query = query.Preload("Streams")
		if args.WithProducers {
			query = query.Preload("Streams.Producers")
			if args.WithWriterClones {
				query = query.Preload("Streams.Producers.WriterClones")
			}
		}
		if args.WithCommandGroups {
			query = query.Preload("Streams.CommandGroups")
		}
	}
	if args.GlobalUUID != nil {
		query = query.Where("global_uuid = ?", *args.GlobalUUID)
	}
	if args.ClientId != nil {
		values := strings.Split(*args.ClientId, ",")
		query = query.Where(fmt.Sprintf(`client_id IN ( '%s' )`, strings.Join(values, "', '")))
	}
	if args.SiteId != nil {
		values := strings.Split(*args.SiteId, ",")
		query = query.Where(fmt.Sprintf(`site_id IN ( '%s' )`, strings.Join(values, "', '")))
	}
	if args.DeviceId != nil {
		values := strings.Split(*args.DeviceId, ",")
		query = query.Where(fmt.Sprintf(`device_id IN ( '%s' )`, strings.Join(values, "', '")))
	}
	if args.SourceUUID != nil {
		values := strings.Split(*args.SourceUUID, ",")
		query = query.Where(fmt.Sprintf(`source_uuid IN ( '%s' )`, strings.Join(values, "', '")))
	}
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	if args.IsRemote != nil {
		if *args.IsRemote {
			query = query.Where("is_remote IS TRUE")
		} else {
			query = query.Where("is_remote IS FALSE")
		}
	}
	return query
}

func (d *GormDatabase) buildFlowNetworkCloneQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithStreamClones {
		query = query.Preload("StreamClones")
		if args.WithTags {
			query = query.Preload("StreamClones.Tags")
		}
		if args.WithConsumers {
			query = query.Preload("StreamClones.Consumers")
			if args.WithTags {
				query = query.Preload("StreamClones.Consumers.Tags")
			}
			if args.WithWriters {
				query = query.Preload("StreamClones.Consumers.Writers")
			}
		}
	}
	if args.GlobalUUID != nil {
		query = query.Where("global_uuid = ?", *args.GlobalUUID)
	}
	if args.SourceUUID != nil {
		query = query.Where("source_uuid = ?", *args.SourceUUID)
	}
	if args.ClientId != nil {
		values := strings.Split(*args.ClientId, ",")
		query = query.Where(fmt.Sprintf(`client_id IN ( '%s' )`, strings.Join(values, "', '")))
	}
	if args.SiteId != nil {
		values := strings.Split(*args.SiteId, ",")
		query = query.Where(fmt.Sprintf(`site_id IN ( '%s' )`, strings.Join(values, "', '")))
	}
	if args.DeviceId != nil {
		values := strings.Split(*args.DeviceId, ",")
		query = query.Where(fmt.Sprintf(`device_id IN ( '%s' )`, strings.Join(values, "', '")))
	}
	if args.UUID != nil {
		query = query.Where("uuid = ?", *args.UUID)
	}
	return query
}

func (d *GormDatabase) buildStreamQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithFlowNetworks {
		query = query.Preload("FlowNetworks")
	}
	if args.WithProducers {
		query = query.Preload("Producers")
		if args.WithWriterClones {
			query = query.Preload("Producers.WriterClones")
		}
	}
	if args.WithCommandGroups {
		query = query.Preload("CommandGroups")
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	return query
}

func (d *GormDatabase) buildStreamCloneQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithConsumers {
		query = query.Preload("Consumers")
		if args.WithWriters {
			query = query.Preload("Consumers.Writers")
		}
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.SourceUUID != nil {
		query = query.Where("source_uuid = ?", *args.SourceUUID)
	}
	if args.FlowNetworkCloneUUID != nil {
		query = query.Where("flow_network_clone_uuid = ?", *args.FlowNetworkCloneUUID)
	}
	return query
}

func (d *GormDatabase) buildConsumerQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithWriters {
		query = query.Preload("Writers")
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.ProducerUUID != nil {
		query = query.Where("producer_uuid = ?", *args.ProducerUUID)
	}
	return query
}

func (d *GormDatabase) buildProducerQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithWriterClones {
		query = query.Preload("WriterClones")
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.StreamUUID != nil {
		query = query.Where("stream_uuid = ?", *args.StreamUUID)
	}
	if args.ProducerThingUUID != nil {
		query = query.Where("producer_thing_uuid = ?", *args.ProducerThingUUID)
	}
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	return query
}

func (d *GormDatabase) buildNetworkQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithDevices {
		query = query.Preload("Devices")
		if args.WithTags {
			query = query.Preload("Devices.Tags")
		}
	}
	if args.WithPoints {
		query = query.Preload("Devices.Points").Preload("Devices.Points.Priority")
		if args.WithTags {
			query = query.Preload("Devices.Points.Tags")
		}
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	return query
}

func (d *GormDatabase) buildDeviceQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithPoints {
		query = query.Preload("Points")
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.WithPriority {
		query = query.Preload("Points.Priority")
	}
	if args.AddressUUID != nil {
		query = query.Where("address_uuid = ?", *args.AddressUUID)
	}
	return query
}

func (d *GormDatabase) buildPointQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithPriority {
		query = query.Preload("Priority")
	}
	if args.WithTags {
		query = query.Preload("Tags")
	}
	if args.AddressUUID != nil {
		query = query.Where("address_uuid = ?", *args.AddressUUID)
	}
	if args.IoNumber != nil {
		query = query.Where("io_number = ?", *args.IoNumber)
	}
	if args.AddressID != nil {
		query = query.Where("address_id = ?", *args.AddressID)
	}
	if args.ObjectType != nil {
		query = query.Where("object_type = ?", *args.ObjectType)
	}
	return query
}

func (d *GormDatabase) buildWriterQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.ConsumerUUID != nil {
		query = query.Where("consumer_uuid = ?", *args.ConsumerUUID)
	}
	if args.WriterThingClass != nil {
		query = query.Where("writer_thing_class = ?", *args.WriterThingClass)
	}
	if args.WriterThingUUID != nil {
		query = query.Where("writer_thing_uuid = ?", *args.WriterThingUUID)
	}
	return query
}

func (d *GormDatabase) buildWriterCloneQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.ProducerUUID != nil {
		query = query.Where("producer_uuid = ?", *args.ProducerUUID)
	}
	if args.WriterThingClass != nil {
		query = query.Where("writer_thing_class = ?", *args.WriterThingClass)
	}
	if args.SourceUUID != nil {
		query = query.Where("source_uuid = ?", *args.SourceUUID)
	}
	return query
}

func (d *GormDatabase) buildTagQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Networks {
		query = query.Preload("Networks")
	}
	if args.WithDevices {
		query = query.Preload("Devices")
	}
	if args.WithPoints {
		query = query.Preload("Points")
	}
	if args.WithStreams {
		query = query.Preload("Streams")
	}
	if args.WithProducers {
		query = query.Preload("Producers")
	}
	if args.WithConsumers {
		query = query.Preload("Consumers")
	}
	return query
}

func (d *GormDatabase) buildProducerHistoryQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.IdGt != nil {
		query = query.Where("Id > ?", args.IdGt)
	}
	if args.TimestampGt != nil {
		query = query.Where("timestamp > datetime(?)", args.TimestampGt)
	}
	if args.TimestampLt != nil {
		query = query.Where("timestamp < datetime(?)", args.TimestampLt)
	}
	if args.Order != "" {
		order := strings.ToUpper(strings.TrimSpace(args.Order))
		if order != "ASC" && order != "DESC" {
			args.Order = "DESC"
		}
		query = query.Order("timestamp " + args.Order)
	}
	return query
}

func (d *GormDatabase) buildHistoryQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.TimestampGt != nil {
		query = query.Where("timestamp > datetime(?)", args.TimestampGt)
	}
	if args.TimestampLt != nil {
		query = query.Where("timestamp < datetime(?)", args.TimestampLt)
	}
	if args.Order != "" {
		order := strings.ToUpper(strings.TrimSpace(args.Order))
		if order != "ASC" && order != "DESC" {
			args.Order = "DESC"
		}
		query = query.Order("timestamp " + args.Order)
	}
	return query
}

func (d *GormDatabase) buildScheduleQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Name != nil {
		query = query.Where("name = ?", *args.Name)
	}
	return query
}
