# CHANGELOG
## [v0.5.5](https://github.com/NubeIO/flow-framework/tree/v0.5.5) (2022-06-10)

- redo of bacnet-server plugin
- fix bug on lora plugin
- reformat nil helpers
- add in sync on flow-networks

## [v0.5.6](https://github.com/NubeIO/flow-framework/tree/v0.5.6) (2022-06-17)

- minor modbus updates

## [v0.5.5](https://github.com/NubeIO/flow-framework/tree/v0.5.5) (2022-06-10)

- redo of bacnet-server plugin
- fix bug on lora plugin
- reformat nil helpers
- add in sync on flow-networks

## [v0.5.4](https://github.com/NubeIO/flow-framework/tree/v0.5.4) (2022-05-24)

- bacnet master fix bug on select object type

## [v0.5.3](https://github.com/NubeIO/flow-framework/tree/v0.5.3) (2022-05-24)

- cascade delete on flow networks
- added native bacnet master plugin

## [v0.5.2](https://github.com/NubeIO/flow-framework/tree/v0.5.2) (2022-05-19)

- Make configurable modbus log level

## [v0.5.1](https://github.com/NubeIO/flow-framework/tree/v0.5.1) (2022-05-19)

- Minor update to polling (safeties on timeouts)
- Added a Postgres plugin

## [v0.5.0](https://github.com/NubeIO/flow-framework/tree/v0.5.0) (2022-05-17)

- fixed modbus point/device/network time settings
- updated pointWrite() for modbus compatibility
- priority array utilities added

## [v0.4.9](https://github.com/NubeIO/flow-framework/tree/v0.4.9) (2022-05-13)

- Fix: isChange checker for COV updates

## [v0.4.8](https://github.com/NubeIO/flow-framework/tree/v0.4.8) (2022-05-12)

- Control COV streams on point value update

## [v0.4.7](https://github.com/NubeIO/flow-framework/tree/v0.4.7) (2022-05-12)

- update to modbus-server plugin

## [v0.4.6](https://github.com/NubeIO/flow-framework/tree/v0.4.6) (2022-05-11)

- added new plugin modbus server
- nullable fallback point field support
- bug fixes on points for priority array
- polling for modbus
- bulk write api on points

## [v0.4.5](https://github.com/NubeIO/flow-framework/tree/v0.4.5) (2022-05-02)

- get bacnet auto mapping working to rubix-io
- fix bug on rubix-io plugin
- fix bug on flow-network sync

## [v0.4.4](https://github.com/NubeIO/flow-framework/tree/v0.4.4) (2022-04-29)

- add new plugin for the rubix-io

## [v0.4.3](https://github.com/NubeIO/flow-framework/tree/v0.4.3) (2022-04-26)

- many updates since last build :)

## [v0.4.2](https://github.com/NubeIO/flow-framework/tree/v0.4.2) (2022-03-23)

- Plugin/testing #417
- Improvement/database plugins #415
- small updates to system plugin #414
- added point math func on point write value #413
- updates to modbus, lora and bacnet #412
- updates to lora and bacnet #411

## [v0.4.1](https://github.com/NubeIO/flow-framework/tree/v0.4.1) (2022-03-17)

- Feature: add present_value on writer #405
- small updates to lora #404
- Improvement/db plugins #403
- Improvement/db plugins #402
- remove on device model AddressUUID unique #401
- Improvement/history #399
- Improvement: nullable address_uuid #397
- fix get networks #396
- Support OR query on client_id, site_id & device_id #395

## [v0.4.0](https://github.com/NubeIO/flow-framework/tree/v0.4.0) (2022-02-28)

- Supporting older schedule deployment

## [v0.3.9](https://github.com/NubeIO/flow-framework/tree/v0.3.9) (2022-02-28)

- Added back GET/PATCH by point's name

## [v0.3.8](https://github.com/NubeIO/flow-framework/tree/v0.3.8) (2022-02-24)

- Improvements/misc #385

## [v0.3.7](https://github.com/NubeIO/flow-framework/tree/v0.3.7) (2022-02-24)

- Improvements/misc #384
- Improvement/stream sync #383
- updates to modbus #382
- Improvement/point value update #381
- Replace modbus lib #379

## [v0.3.6](https://github.com/NubeIO/flow-framework/tree/v0.3.6) (2022-02-22)

- Merge pull request #358 from NubeIO/bac-master-1
- Fix: edit flow-network issue
- Fix: issue on P2P openvpn connected devices
  - SyncFlowNetwork: Post "10.8.1.1:1616/ff/api/sync/flow_network";: EOF
- updates to modbus plugin
- Add writers write, read, sync action
- Fix: writer action support for old deployments

## [v0.3.5](https://github.com/NubeIO/flow-framework/tree/v0.3.5) (2022-02-16)

- Improvement: support multiple producers on a single point
- Improvement on producer history
- Writer action updates its own side thing
- Make schedule write actions working
- Fix: FlowNetwork creation issue for HTTP
- Make FlowNetwork update working
- Fix: float pointer values comparison for COV

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
