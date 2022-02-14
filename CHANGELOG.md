# CHANGELOG
## [v0.3.4](https://github.com/NubeIO/flow-framework/tree/v0.3.4) (2022-02-10)
- Sync value on COV of point
- Fix: point.present_value comparison issue
- Fix: database lock issue

## [v0.3.3](https://github.com/NubeIO/flow-framework/tree/v0.3.3) (2022-02-10)
- Fix: schedule

## [v0.3.2](https://github.com/NubeIO/flow-framework/tree/v0.3.2) (2022-02-10)
- small bug fix to schedule small fix to stop sch crashing #352

## [v0.3.1](https://github.com/NubeIO/flow-framework/tree/v0.3.1) (2022-02-10)
- small bug fix to schedule small fix to stop sch crashing #351

## [v0.3.0](https://github.com/NubeIO/flow-framework/tree/v0.3.0) (2022-02-10)
- Add schedule writer POC
- Add with_priority option on device query builder
- Improvements on gorm migration #347
- added bacnetmaster plugin #346
- make sure if pnt is same the addrID is not same #345
- clean up of bacnet-server #343
- Update/schedule checker to new schedule json schema #342
- small fixes to modbus #337
- Improvements on schedule #334
- Update to mqtt broker plugin
- Add schedule config
- Sync schedules values on the producer side
- Remove CreateWriterWizard
- Datastore is nil for updating writers (datastore update is only from write actions)
- Add sync on patch (#326)
- Marc/edge28 plugin scaling (#314)
- Merge pull request #325 from NubeIO/improvement/return-appropriate-status-code

## [v0.2.2](https://github.com/NubeIO/flow-framework/tree/v0.2.2) (2021-12-23)
- Improvement on schedule APIs
- Add scheduler for refreshing token

## [v0.2.1](https://github.com/NubeIO/flow-framework/tree/v0.2.1) (2021-12-18)
- Update: update to schedules

## [v0.2.0](https://github.com/NubeIO/flow-framework/tree/v0.2.0) (2021-12-18)
- Update: update to schedules

## [v0.1.9](https://github.com/NubeIO/flow-framework/tree/v0.1.9) (2021-12-16)
- Remove: rubix plugins

## [v0.1.8](https://github.com/NubeIO/flow-framework/tree/v0.1.71) (2021-11-23)
- Fix: get config on sessions

## [v0.1.71](https://github.com/NubeIO/flow-framework/tree/v0.1.71) (2021-11-22)
- made rubix a network #304
- Breaking issue fix on modubs polling #302 
- Improvements/misc #301
- Improvement/misc #300
- Add history influx log #299
- Fix: NubeIO vs NubeDev package #298
- Feature/schema api #297
- Update go.mod #296
- updated the bash script #295
- Improvement/misc #294
- Close DB connection #293
- sample api helper #290

## [v0.1.6](https://github.com/NubeIO/flow-framework/tree/v0.1.6) (2021-11-04)
- add APIs for proxying fn, fnc
- serial port fix (#288)

## [v0.1.5](https://github.com/NubeIO/flow-framework/tree/v0.1.5) (2021-11-01)
- nubeio-rubix-lib-helpers-go version upgrade to v0.1.2

## [v0.1.4](https://github.com/NubeIO/flow-framework/tree/v0.1.4) (2021-10-25)
- rubix plugin build fix

## [v0.1.3](https://github.com/NubeIO/flow-framework/tree/v0.1.3) (2021-10-25)
- added rubix-service api
- fix up on schedules
- improvement on writer & writer_clone args query

## [v0.1.2](https://github.com/NubeIO/flow-framework/tree/v0.1.2) (2021-10-21)
- added flow network mqtt api

## [v0.1.1](https://github.com/NubeIO/flow-framework/tree/v0.1.1) (2021-10-19)
- added api for milo db
- clean up of bacnetserver plugin
- added system and time api
- added schedules api

## [v0.1.0](https://github.com/NubeIO/flow-framework/tree/v0.1.0) (2021-10-12)
- fix issues on droplet motion
- added writer-action as thingClass schedule

## [v0.0.9](https://github.com/NubeIO/flow-framework/tree/v0.0.9) (2021-10-11)
- updates to lora and modbus plugins
- added edge-28 plugin

## [v0.0.8](https://github.com/NubeIO/flow-framework/tree/v0.0.8) (2021-10-08)
- updates to lora and modbus plugins

## [v0.0.7](https://github.com/NubeIO/flow-framework/tree/v0.0.7) (2021-10-05)
- updates to lora and modbus plugin

## [v0.0.6](https://github.com/NubeIO/flow-framework/tree/v0.0.6) (2021-10-05)
- fix bug on action write

## [v0.0.5](https://github.com/NubeIO/flow-framework/tree/v0.0.5) (2021-10-05)
- added point calc's, units, and eval
- clean up on lora and modbus plugins

## [v0.0.4](https://github.com/NubeIO/flow-framework/tree/v0.0.4) (2021-10-02)
- added git plugin and updates to modbus

## [v0.0.3](https://github.com/NubeIO/flow-framework/tree/v0.0.3) (2021-10-01)
- make that artifacts working for armv7

## [v0.0.2](https://github.com/NubeIO/flow-framework/tree/v0.0.2) (2021-09-29)
- include plugins on artifacts

## [v0.0.1](https://github.com/NubeIO/flow-framework/tree/v0.0.1) (2021-08-26)
- first initial release
