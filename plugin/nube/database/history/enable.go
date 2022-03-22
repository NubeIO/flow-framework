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
	_, _ = cron.Every(i.config.Job.Frequency).Tag("SyncHistory").Do(i.syncHistory)
	cron.StartAsync()
	return nil
}

func (i *Instance) Disable() error {
	i.enabled = false
	cron.Clear()
	return nil
}
