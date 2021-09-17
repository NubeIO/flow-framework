

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
