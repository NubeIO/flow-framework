package pollqueue

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	log "github.com/sirupsen/logrus"
)

func (pm *NetworkPollManager) pollQueueDebugMsg(args ...interface{}) {
	if nstring.InEqualIgnoreCase(pm.Config.LogLevel, "DEBUG") {
		prefix := "Modbus Poll Queue: "
		log.Info(prefix, args)
	}
}

func (pm *NetworkPollManager) pollQueueErrorMsg(args ...interface{}) {
	prefix := "Modbus Poll Queue: "
	log.Error(prefix, args)
}

func (pm *NetworkPollManager) PrintPollQueuePointUUIDs() {
	if nstring.InEqualIgnoreCase(pm.Config.LogLevel, "DEBUG") { //Added here to disable debug processes when not using logging
		printString := "\n\n"
		printString += fmt.Sprint("PrintPollQueuePointUUIDs: (NOTE: THE CURRENT PollPoint HAS ALREADY BEEN REMOVED FROM THE QUEUES AT THIS POINT!!\nTOTAL COUNT = ", pm.PollQueue.PriorityQueue.Len()+pm.PollQueue.StandbyPollingPoints.Len()+pm.PollQueue.OutstandingPollingPoints.Len(), "\n")
		printString += fmt.Sprint("NextPollPoint: ")
		if pm.PluginQueueUnloader != nil {
			printString += fmt.Sprintf("%+v\n", pm.PluginQueueUnloader.NextPollPoint)
		}
		printString += fmt.Sprint("PollQueue: COUNT = ", pm.PollQueue.PriorityQueue.Len(), ": ")
		for _, pp := range pm.PollQueue.PriorityQueue.PriorityQueue {
			printString += fmt.Sprint(pp.FFPointUUID, " - ", pp.PollPriority, "; ")
		}
		printString += fmt.Sprint("", "\n")
		printString += fmt.Sprint("StandbyPollingPoints COUNT = ", pm.PollQueue.StandbyPollingPoints.Len(), ": ")
		for _, pp := range pm.PollQueue.StandbyPollingPoints.PriorityQueue {
			printString += fmt.Sprint(pp.FFPointUUID, " - ", pp.PollPriority, ", repoll:", pp.RepollTimer != nil, "; ")
		}
		printString += fmt.Sprint("\n")
		printString += fmt.Sprint("CurrentlyPollingPoints COUNT = ", pm.PollQueue.OutstandingPollingPoints.Len(), ": ")
		for _, pp := range pm.PollQueue.OutstandingPollingPoints.PriorityQueue {
			printString += fmt.Sprint(pp.FFPointUUID, " - ", pp.PollPriority, ", repoll:", pp.RepollTimer != nil, "; ")
		}
		printString += fmt.Sprint("\n")
		pm.pollQueueDebugMsg(printString)
	}
}

func (pm *NetworkPollManager) PrintPointDebugInfo(pnt *model.Point) {
	if nstring.InEqualIgnoreCase(pm.Config.LogLevel, "DEBUG") { //Added here to disable debug processes when not using logging
		printString := "\n\n"
		if pnt != nil {
			printString += fmt.Sprint("Point: ", pnt.UUID, " ", pnt.Name, "\n")
			printString += fmt.Sprint("WriteMode: ", pnt.WriteMode, "\n")
			if pnt.WritePollRequired != nil {
				printString += fmt.Sprint("WritePollRequired: ", *pnt.WritePollRequired, "\n")
			}
			if pnt.ReadPollRequired != nil {
				printString += fmt.Sprint("ReadPollRequired: ", *pnt.ReadPollRequired, "\n")
			}
			if pnt.WriteValue == nil {
				printString += fmt.Sprint("WriteValue: nil", "\n")
			} else {
				printString += fmt.Sprint("WriteValue: ", *pnt.WriteValue, "\n")
			}
			if pnt.OriginalValue == nil {
				printString += fmt.Sprint("OriginalValue: nil", "\n")
			} else {
				printString += fmt.Sprint("OriginalValue: ", *pnt.OriginalValue, "\n")
			}
			if pnt.PresentValue == nil {
				printString += fmt.Sprint("PresentValue: nil", "\n")
			} else {
				printString += fmt.Sprint("PresentValue: ", *pnt.PresentValue, "\n")
			}
			if pnt.CurrentPriority == nil {
				printString += fmt.Sprint("CurrentPriority: nil", "\n")
			} else {
				printString += fmt.Sprint("CurrentPriority: ", *pnt.CurrentPriority, "\n")
			}
			if pnt.Priority != nil {
				if pnt.Priority.P1 != nil {
					printString += fmt.Sprint("_1: ", *pnt.Priority.P1, "\n")
				}
				if pnt.Priority.P2 != nil {
					printString += fmt.Sprint("_2: ", *pnt.Priority.P2, "\n")
				}
				if pnt.Priority.P3 != nil {
					printString += fmt.Sprint("_3: ", *pnt.Priority.P3, "\n")
				}
				if pnt.Priority.P4 != nil {
					printString += fmt.Sprint("_4: ", *pnt.Priority.P4, "\n")
				}
				if pnt.Priority.P5 != nil {
					printString += fmt.Sprint("_5: ", *pnt.Priority.P5, "\n")
				}
				if pnt.Priority.P6 != nil {
					printString += fmt.Sprint("_6: ", *pnt.Priority.P6, "\n")
				}
				if pnt.Priority.P7 != nil {
					printString += fmt.Sprint("_7: ", *pnt.Priority.P7, "\n")
				}
				if pnt.Priority.P8 != nil {
					printString += fmt.Sprint("_8: ", *pnt.Priority.P8, "\n")
				}
				if pnt.Priority.P9 != nil {
					printString += fmt.Sprint("_9: ", *pnt.Priority.P9, "\n")
				}
				if pnt.Priority.P10 != nil {
					printString += fmt.Sprint("_10: ", *pnt.Priority.P10, "\n")
				}
				if pnt.Priority.P11 != nil {
					printString += fmt.Sprint("_11: ", *pnt.Priority.P11, "\n")
				}
				if pnt.Priority.P12 != nil {
					printString += fmt.Sprint("_12: ", *pnt.Priority.P12, "\n")
				}
				if pnt.Priority.P13 != nil {
					printString += fmt.Sprint("_13: ", *pnt.Priority.P13, "\n")
				}
				if pnt.Priority.P14 != nil {
					printString += fmt.Sprint("_14: ", *pnt.Priority.P14, "\n")
				}
				if pnt.Priority.P15 != nil {
					printString += fmt.Sprint("_15: ", *pnt.Priority.P15, "\n")
				}
				if pnt.Priority.P16 != nil {
					printString += fmt.Sprint("_16: ", *pnt.Priority.P16, "\n")
				}
			}
			pm.pollQueueDebugMsg(printString)
			return
		}
		pm.pollQueueDebugMsg("ERROR: INVALID POINT")
	}
}

func (pm *NetworkPollManager) PrintPollingPointDebugInfo(pp *PollingPoint) {
	if pp != nil {
		pm.pollQueueDebugMsg(fmt.Sprintf("ModbusPolling() pp %+v", pp))
	}
}

func (pm *NetworkPollManager) PrintPollQueueStatistics() {
	if nstring.InEqualIgnoreCase(pm.Config.LogLevel, "DEBUG") { //Added here to disable debug processes when not using logging

		pm.PrintPollQueuePointUUIDs()

		printString := "\n\n"
		printString += fmt.Sprint("PrintPollQueueStatistics: \n")
		printString += fmt.Sprint("MaxPollExecuteTimeSecs: ", pm.MaxPollExecuteTimeSecs, "\n")
		printString += fmt.Sprint("AveragePollExecuteTimeSecs: ", pm.AveragePollExecuteTimeSecs, "\n")
		printString += fmt.Sprint("MinPollExecuteTimeSecs: ", pm.MinPollExecuteTimeSecs, "\n")
		printString += fmt.Sprint("TotalPollQueueLength: ", pm.TotalPollQueueLength, "\n")
		printString += fmt.Sprint("TotalStandbyPointsLength: ", pm.TotalStandbyPointsLength, "\n")
		printString += fmt.Sprint("TotalPointsOutForPolling: ", pm.TotalPointsOutForPolling, "\n")
		printString += fmt.Sprint("ASAPPriorityPollQueueLength: ", pm.ASAPPriorityPollQueueLength, "\n")
		printString += fmt.Sprint("HighPriorityPollQueueLength: ", pm.HighPriorityPollQueueLength, "\n")
		printString += fmt.Sprint("NormalPriorityPollQueueLength: ", pm.NormalPriorityPollQueueLength, "\n")
		printString += fmt.Sprint("LowPriorityPollQueueLength: ", pm.LowPriorityPollQueueLength, "\n")
		printString += fmt.Sprint("ASAPPriorityAveragePollTime: ", pm.ASAPPriorityAveragePollTime, "\n")
		printString += fmt.Sprint("HighPriorityAveragePollTime: ", pm.HighPriorityAveragePollTime, "\n")
		printString += fmt.Sprint("NormalPriorityAveragePollTime: ", pm.NormalPriorityAveragePollTime, "\n")
		printString += fmt.Sprint("LowPriorityAveragePollTime: ", pm.LowPriorityAveragePollTime, "\n")
		printString += fmt.Sprint("TotalPollCount: ", pm.TotalPollCount, "\n")
		printString += fmt.Sprint("ASAPPriorityPollCount: ", pm.ASAPPriorityPollCount, "\n")
		printString += fmt.Sprint("HighPriorityPollCount: ", pm.HighPriorityPollCount, "\n")
		printString += fmt.Sprint("NormalPriorityPollCount: ", pm.NormalPriorityPollCount, "\n")
		printString += fmt.Sprint("LowPriorityPollCount: ", pm.LowPriorityPollCount, "\n")
		printString += fmt.Sprint("ASAPPriorityPollCountForAvg: ", pm.ASAPPriorityPollCountForAvg, "\n")
		printString += fmt.Sprint("HighPriorityPollCountForAvg: ", pm.HighPriorityPollCountForAvg, "\n")
		printString += fmt.Sprint("NormalPriorityPollCountForAvg: ", pm.NormalPriorityPollCountForAvg, "\n")
		printString += fmt.Sprint("LowPriorityPollCountForAvg: ", pm.LowPriorityPollCountForAvg, "\n")
		printString += fmt.Sprint("ASAPPriorityMaxCycleTime: ", pm.ASAPPriorityMaxCycleTime, "\n")
		printString += fmt.Sprint("HighPriorityMaxCycleTime: ", pm.HighPriorityMaxCycleTime, "\n")
		printString += fmt.Sprint("NormalPriorityMaxCycleTime: ", pm.NormalPriorityMaxCycleTime, "\n")
		printString += fmt.Sprint("LowPriorityMaxCycleTime: ", pm.LowPriorityMaxCycleTime, "\n")
		printString += fmt.Sprint("ASAPPriorityLockupAlert: ", pm.ASAPPriorityLockupAlert, "\n")
		printString += fmt.Sprint("HighPriorityLockupAlert: ", pm.HighPriorityLockupAlert, "\n")
		printString += fmt.Sprint("NormalPriorityLockupAlert: ", pm.NormalPriorityLockupAlert, "\n")
		printString += fmt.Sprint("LowPriorityLockupAlert: ", pm.LowPriorityLockupAlert, "\n")
		printString += fmt.Sprint("PollingStartTimeUnix: ", pm.PollingStartTimeUnix, "\n")
		printString += fmt.Sprint("BusyTime: ", pm.BusyTime, "% \n")
		printString += fmt.Sprint("EnabledTime: ", pm.EnabledTime, "% \n")
		printString += fmt.Sprint("PortUnavailableTime: ", pm.PortUnavailableTime, " \n")
		printString += fmt.Sprint("\n")
		pm.pollQueueDebugMsg(printString)
	}
}
