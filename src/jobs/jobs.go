package jobs

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/eventbus"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/src/dbhandler"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"time"
)

type Jobs struct {
	db      dbhandler.Handler
	Enabled bool
}

var cron *gocron.Scheduler
var bus eventbus.BusService
var enabled bool

//GetJobService will return the instance of the job service
func GetJobService() (*gocron.Scheduler, bool) {
	if enabled {
		return cron, true
	}
	return cron, false
}

func (j *Jobs) InitCron() {
	bus = eventbus.NewService(eventbus.GetBus())
	cron = gocron.NewScheduler(time.UTC)
	cron.StartAsync()
	j.Enabled = true
	enabled = j.Enabled
}

func (j *Jobs) task(mp string, uuid string) {
	t := fmt.Sprintf("%s.%s.%s", eventbus.JobTrigger, mp, uuid)
	log.Info(fmt.Sprintf("Triggering a job: %s", t))
	bus.RegisterTopic(t)
	err := bus.Emit(eventbus.CTX(), t, "MEG OVER BUS")
	if err != nil {
		//TODO FIX ERROR
	}
}

func (j *Jobs) JobAdd(body *model.Job) error {
	pc, err := j.db.GetPlugin(body.PluginConfId)
	if err != nil {
		return err
	}
	_, err = cron.Every(body.Frequency).Tag(body.UUID).Do(j.task, pc.ModulePath, body.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jobs) jobRemover(uuid string) error {
	err := cron.RemoveByTag(uuid)
	if err != nil {
		return errors.New("error on remove job")
	}
	return nil
}
