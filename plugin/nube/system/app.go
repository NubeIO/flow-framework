package main

import (
	"github.com/NubeDev/flow-framework/src/jobs"
	log "github.com/sirupsen/logrus"
)

func (i *Instance) schedule() {
	j, ok := jobs.GetJobService()
	if ok {
		_, err := j.Every(30).Second().Do(i.run)
		if err != nil {
			log.Infof("system-plugin-schedule: error on create job %v\n", err)
		}
	}
}
