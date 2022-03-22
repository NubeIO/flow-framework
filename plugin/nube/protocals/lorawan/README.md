# scope

- a network in `ff` with be a `cs` gateway
- a device will be a device
- a point will be a sensor value from a device
- the org id will be default be set to 1 and will be set in the config file
- MQTT lister `application/+/device/+/rx`

### mqtt payload

```
{
  "topic": "application/1/device/a81758fffe064d01/rx",
  "payload": {
    "applicationID": "1",
    "applicationName": "default-app",
    "deviceName": "LRWAN_WATER1",
    "devEUI": "a81758fffe064d01",
    "rxInfo": [
      {
        "gatewayID": "0000000000000000",
        "uplinkID": "c8ee256c-385a-44ea-bd26-5c708fe24a80",
        "name": "local-rak-gateway",
        "rssi": -50,
        "loRaSNR": 9.2,
        "location": {
          "latitude": 0,
          "longitude": 0,
          "altitude": 0
        }
      }
    ],
    "txInfo": {
      "frequency": 916400000,
      "dr": 0
    },
    "adr": true,
    "fCnt": 2175,
    "fPort": 5,
    "data": "CwAAAAsXAAAAAA==",
    "object": {
      "pulseAbs": 11,
      "pulseAbs2": 0
    }
  },
  "qos": 0,
  "retain": false,
  "_msgid": "d0a43c82291aabb9"
}
```

### example

```
curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' --header 'Grpc-Metadata-Authorization: Bearer <API TOKEN>' -d '{ \ 
   "deviceQueueItem": { \ 
     "confirmed": false, \ 
     "data": "AQID", \ 
     "fPort": 10 \ 
   } \ 
 }' 'http://localhost:8080/api/devices/0101010101010101/queue'
```


### swagger
```
http://localhost:8080/api
```
