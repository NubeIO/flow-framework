# Scope
- Syncs data to InfluxDB2
- Flow: table.histories > InfluxDB2 
    - Organization: nube-org
    - Bucket: nube-bucket
- It uses config file for InfluxDB connection & job to sync that value
  - Use Token for credential; `influx auth list`
- It uses Job for starting scheduler task


### How to get default config
- Save empty `YAML` file, and it will generate you the default config file

### How to delete influxDB data
- `influx delete --org nube-org --bucket nube-bucket --start 2020-03-01T00:00:00Z --stop 2023-11-14T00:00:00Z`