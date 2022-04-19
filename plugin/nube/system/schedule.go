package main

import (
	"github.com/NubeIO/flow-framework/src/schedule"
	"github.com/NubeIO/flow-framework/src/utilstime"
	"github.com/NubeIO/flow-framework/utils"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (i *Instance) runSchedule() {
	schedules, err := i.db.GetSchedules()
	if err != nil {
		return
	}
	for _, sch := range schedules {
		ScheduleJSON, err := schedule.DecodeSchedule(sch.Schedule)
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
			return
		}

		scheduleNameToCheck := "ALL" //TODO: we need a way to specify the schedule name that is being checked for.

		timezone := ScheduleJSON.Config.TimeZone
		_, err = time.LoadLocation(timezone)
		if err != nil || timezone == "" { // If timezone field is not assigned or invalid, get timezone from System Time
			log.Error("system-plugin-schedule: CheckWeeklyScheduleCollection(): invalid schedule timezone. checking with system time.")
			systemTimezone := strings.Split((*utilstime.SystemTime()).HardwareClock.Timezone, " ")[0]
			//fmt.Println("systemTimezone 2: ", systemTimezone)
			if systemTimezone == "" {
				zone, _ := utilstime.GetHardwareTZ()
				timezone = zone
			} else {
				timezone = systemTimezone
			}
		}

		//CHECK WEEKLY SCHEDULES
		weeklyResult, err := schedule.WeeklyCheck(ScheduleJSON.Schedules.Weekly, scheduleNameToCheck, timezone) //This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
		}
		// CHECK EVENT SCHEDULES
		eventResult, err := schedule.EventCheck(ScheduleJSON.Schedules.Events, scheduleNameToCheck, timezone) //This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on EventCheck %v\n", err)
		}
		//Combine Event and Weekly schedule results.
		weeklyAndEventResult, err := schedule.CombineScheduleCheckerResults(weeklyResult, eventResult)
		// CHECK EXCEPTION SCHEDULES
		exceptionResult, err := schedule.ExceptionCheck(ScheduleJSON.Schedules.Exceptions, scheduleNameToCheck, timezone) //This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on ExceptionCheck %v\n", err)
		}
		if exceptionResult.CheckIfEmpty() {
			log.Println("system-plugin-schedule: Exception schedule is empty")
		}

		finalResult, err := schedule.ApplyExceptionSchedule(weeklyAndEventResult, exceptionResult) //This applies the exception schedule to mask the combined weekly and event schedules.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on ApplyExceptionSchedule %v\n", err)
		}

		log.Printf("system-plugin-schedule: finalResult: %+v\n", finalResult.IsActive)

		if sch != nil {
			i.store.Set(sch.Name, finalResult, -1)
			s := new(model.Schedule)
			if finalResult.IsActive {
				s.IsActive = utils.NewTrue()
			} else {
				s.IsActive = utils.NewFalse()
			}
			if utils.IsTrue(s.IsActive) != utils.IsTrue(sch.IsActive) {
				log.Printf("system-plugin-schedule: UPDATE SCHEDULE IN DB %v\n", sch.Name)
				_, err = i.db.UpdateSchedule(sch.UUID, s)
				if err != nil {
					log.Errorf("system-plugin-schedule: issue on UpdateSchedule %v\n", sch.UUID)
				}
			}
		}

	}
}
