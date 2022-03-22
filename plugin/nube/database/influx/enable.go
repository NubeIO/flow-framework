package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (i *Instance) Enable() error {
	i.enabled = true
	i.setUUID()
	i.BusServ()
	cron = gocron.NewScheduler(time.UTC)
	influxDetails := i.initializeInfluxSettings()
	_, _ = cron.Every(i.config.Job.Frequency).Tag("SyncInflux").Do(i.syncInflux, influxDetails)
	cron.StartAsync()
	return nil
}

func (i *Instance) Disable() error {
	i.enabled = false
	cron.Clear()
	for _, influxConnectionInstance := range influxConnectionInstances {
		influxConnectionInstance.client.Close()
	}
	return nil
}
