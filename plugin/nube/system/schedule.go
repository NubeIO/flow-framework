package main

/*
func (i *Instance) run() {
	class, err := i.db.GetWriters(api.Args{WriterThingClass: utils.NewStringAddress("schedule")})
	if err != nil {
		return
	}
	for _, v := range class {
		decodeSchedule, err := schedule.DecodeSchedule(v.DataStore)
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
		}
		scheduleNameToCheck := "HVAC" //TODO: we need a way to specify the schedule name that is being checked for.

		// CHECK WEEKLY SCHEDULES
		//result, err := schedule.WeeklyCheck(decodeSchedule.Weekly, "ANY")  //This will check for any active schedules with any name
		weeklyResult, err := schedule.WeeklyCheck(decodeSchedule.Weekly, scheduleNameToCheck) //This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
		}

		// CHECK EVENT SCHEDULES
		//eventResult, err := schedule.EventCheck(decodeSchedule.Events, "ANY")  //This will check for any active schedules with any name
		eventResult, err := schedule.EventCheck(decodeSchedule.Events, scheduleNameToCheck) //This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on EventCheck %v\n", err)
		}

		//Combine Event and Weekly schedule results.
		weeklyAndEventResult, err := schedule.CombineScheduleCheckerResults(weeklyResult, eventResult)

		// CHECK EXCEPTION SCHEDULES
		//exceptionResult, err := schedule.ExceptionCheck(decodeSchedule.Exceptions, "ANY")  //This will check for any active schedules with any name
		exceptionResult, err := schedule.ExceptionCheck(decodeSchedule.Exceptions, scheduleNameToCheck) //This will check for any active schedules with defined name.
		if err != nil {
			log.Errorf("system-plugin-schedule: issue on ExceptionCheck %v\n", err)
		}

		for schKey, week := range decodeSchedule.Weekly {
			result, err := schedule.WeeklyCheck(decodeSchedule.Weekly, week.Name)
			if err != nil {
				log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
			}
			if sch.IsActive {
				i.store.Set(week.Name, sch, -1)
			}
			log.Infof("system-plugin-schedule: schedule schKey %v\n", schKey)
			log.Infof("system-plugin-schedule: schedule Name %v\n", week.Name)
			log.Infof("system-plugin-schedule: schedule NextStart %v\n", time.Unix(sch.NextStart, 0))
			log.Infof("system-plugin-schedule: schedule NextStop %v\n", time.Unix(sch.NextStop, 0))
			log.Infof("system-plugin-schedule: schedule is IsActive %v\n", sch.IsActive)
			log.Infof("system-plugin-schedule: schedule Payload %v\n", sch.Payload)
		}

	}

}

*/

//func (i *Instance) runScheduleAPI() {
//	schedules, err := i.db.GetSchedules()
//	if err != nil {
//		log.Infof("system-plugin-schedule: db get scheudles %v\n", err)
//	}
//	for _, s := range schedules {
//		decodeSchedule, err := schedule.DecodeSchedule(s.Schedules)
//		if err != nil {
//			log.Errorf("system-plugin-schedule: issue on DecodeSchedule %v\n", err)
//		}
//		for schKey, week := range decodeSchedule.Weekly {
//			sch, err := schedule.WeeklyCheck(decodeSchedule.Weekly, week.Name)
//			if err != nil {
//				log.Errorf("system-plugin-schedule: issue on WeeklyCheck %v\n", err)
//			}
//
//			if sch.IsActive {
//				i.store.Set(week.Name, sch, -1)
//			}
//			log.Infof("system-plugin-schedule: schedule schKey %v\n", schKey)
//			log.Infof("system-plugin-schedule: schedule Name %v\n", week.Name)
//			log.Infof("system-plugin-schedule: schedule NextStart %v\n", time.Unix(sch.NextStart, 0))
//			log.Infof("system-plugin-schedule: schedule NextStop %v\n", time.Unix(sch.NextStop, 0))
//			log.Infof("system-plugin-schedule: schedule is IsActive %v\n", sch.IsActive)
//			log.Infof("system-plugin-schedule: schedule Payload %v\n", sch.Payload)
//		}
//	}
//}
