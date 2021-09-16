package database

import (
	"github.com/NubeDev/flow-framework/api"
	"gorm.io/gorm"
)

func (d *GormDatabase) buildFlowNetworkQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Streams {
		query = query.Preload("Streams")
	}
	return query
}

func (d *GormDatabase) buildStreamQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.FlowNetworks {
		query = query.Preload("FlowNetworks")
	}
	if args.Producers {
		query = query.Preload("Producers")
		if args.Writers {
			query = query.Preload("Producers.WriterClones")
		}
	}
	if args.Consumers {
		query = query.Preload("Consumers")
		if args.Writers {
			query = query.Preload("Consumers.Writers")
		}
	}
	if args.CommandGroups {
		query = query.Preload("CommandGroups")
	}
	return query
}

func (d *GormDatabase) buildConsumerQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Writers {
		query = query.Preload("Writers")
	}
	return query
}

func (d *GormDatabase) buildProducerQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Writers {
		query = query.Preload("WriterClones")
	}
	return query
}

func (d *GormDatabase) buildNetworkQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Devices {
		query = query.Preload("Devices")
	}
	if args.Points {
		query = query.Preload("Devices.Points")
	}
	if args.IpConnection {
		query = query.Preload("IpConnection")
	}
	if args.SerialConnection {
		query = query.Preload("SerialConnection")
	}
	return query
}

func (d *GormDatabase) buildDeviceQuery(args api.Args) *gorm.DB {
	query := d.DB
	if args.Points {
		query = query.Preload("Points")
	}
	return query
}
