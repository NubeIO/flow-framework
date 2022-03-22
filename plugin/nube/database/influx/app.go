package main

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/model"
	log "github.com/sirupsen/logrus"
	"math"
)

type InfluxDetail struct {
	InfluxSetting *InfluxSetting
	MaxId         int
	Records       int
	IsError       bool
}

func (i *Instance) initializeInfluxSettings() []*InfluxSetting {
	var influxSettings []*InfluxSetting
	influxConnections := i.config.Influx
	for _, influx := range influxConnections {
		influxSetting := new(InfluxSetting)
		schema := "http"
		if influx.Port == 443 {
			schema = "https"
		}
		influxSetting.ServerURL = fmt.Sprintf("%s://%s:%d", schema, influx.Host, influx.Port)
		if influx.Token == nil {
			log.Warn("Token is null, please update it")
			continue
		}
		influxSetting.AuthToken = *influx.Token
		influxSetting.Org = influx.Org
		influxSetting.Bucket = influx.Bucket
		influxSetting.Measurement = influx.Measurement
		influxSettings = append(influxSettings, influxSetting)
	}
	return influxSettings
}

func (i *Instance) syncInflux(influxSettings []*InfluxSetting) (bool, error) {
	log.Info("InfluxDB sync has is been called...")
	if len(influxSettings) == 0 {
		err := "InfluxDB sync failure: no any valid InfluxDB connection with not NULL token"
		log.Warn(err)
		return false, errors.New(err)
	}

	leastLastSyncId := math.MinInt16
	var influxDetails []InfluxDetail
	allError := true
	for _, influxSetting := range influxSettings {
		lastSyncId, isError := influxSetting.GetLastSyncId()
		influxDetail := InfluxDetail{
			InfluxSetting: influxSetting,
			MaxId:         lastSyncId,
			Records:       0,
			IsError:       isError,
		}
		if !isError {
			allError = false
		}
		if leastLastSyncId > lastSyncId && !isError {
			leastLastSyncId = lastSyncId
		}
		influxDetails = append(influxDetails, influxDetail)
	}

	if allError {
		err := "no connections are valid"
		log.Warn(err)
		return false, errors.New(err)
	}
	histories, err := i.db.GetHistoriesForSync(leastLastSyncId)
	if err != nil {
		return false, err
	}

	producerUuid := ""
	var historyTags []*model.HistoryInfluxTag
	for _, history := range histories {
		if producerUuid != history.UUID {
			producerUuid = history.UUID
			historyTags, err = i.db.GetHistoryInfluxTags(producerUuid)
			if err != nil || len(historyTags) == 0 {
				log.Warn(fmt.Sprintf("We unable to get the producer_uuid = %s details!", producerUuid))
				continue
			}
		}
		for _, historyTag := range historyTags {
			tags := tagsHistory(historyTag)
			fields := fieldsHistory(history)
			for i, influxDetail := range influxDetails {
				if influxDetail.MaxId < history.ID {
					influxDetail.InfluxSetting.WriteHistories(tags, fields, history.Timestamp)
					influxDetails[i].Records += 1 // directly updating to reflect value
				}
			}
		}
	}

	// forcing to push bulk writes
	for _, influxDetail := range influxDetails {
		influxDetail.InfluxSetting.getInfluxConnectionInstance().writeAPI.Flush()
		if influxDetail.Records > 0 {
			log.Info(fmt.Sprintf("Stored %v rows on %v", influxDetail.Records, path))
		} else {
			log.Info("Nothing to store, no new records")
		}
	}
	return true, nil
}
