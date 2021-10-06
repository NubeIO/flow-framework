package database

import (
	"github.com/NubeDev/flow-framework/api"
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
		query = query.Where("client_id = ?", *args.ClientId)
	}
	if args.SiteId != nil {
		query = query.Where("site_id = ?", *args.SiteId)
	}
	if args.DeviceId != nil {
		query = query.Where("device_id = ?", *args.DeviceId)
	}
	return query
}

func (d *GormDatabase) buildFlowNetworkCloneQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithStreamClones {
		query = query.Preload("StreamClones")
		if args.WithConsumers {
			query = query.Preload("StreamClones.Consumers")
			if args.WithWriters {
				query = query.Preload("StreamClones.Consumers.Writers")
			}
		}
	}
	if args.GlobalUUID != nil {
		query = query.Where("global_uuid = ?", *args.GlobalUUID)
	}
	if args.ClientId != nil {
		query = query.Where("client_id = ?", *args.ClientId)
	}
	if args.SiteId != nil {
		query = query.Where("site_id = ?", *args.SiteId)
	}
	if args.DeviceId != nil {
		query = query.Where("device_id = ?", *args.DeviceId)
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
	return query
}

func (d *GormDatabase) buildNetworkQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.WithDevices {
		query = query.Preload("Devices")
	}
	if args.WithPoints {
		query = query.Preload("Devices.Points").Preload("Devices.Points.Priority")
	}
	if args.WithIpConnection {
		query = query.Preload("IpConnection")
	}
	if args.WithSerialConnection {
		query = query.Preload("SerialConnection")
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
