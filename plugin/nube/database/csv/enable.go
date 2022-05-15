package main

import (
	"github.com/go-co-op/gocron"
	"time"
)

var cron *gocron.Scheduler

func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	cron = gocron.NewScheduler(time.UTC)
	_, _ = cron.Every(inst.config.Job.Frequency).Tag("ExportCOV").Do(inst.ExportCSV)
	cron.StartAsync()
	return nil
}

func (inst *Instance) Disable() error {
	inst.enabled = false
	cron.Clear()
	return nil
}
