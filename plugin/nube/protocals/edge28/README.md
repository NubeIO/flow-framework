## scope

this plugin is for use api to edge-28 controller

## api's used

ping

``` 
http://192.168.15.101:5000/
```

write to UOs UOs 0 = 12vdc and 100 = 0vdc (Yes its backwards) <io_num>// uo/uo1/22/16 the priority (pri) is not
supported yet but it's there for future use if needed

``` 
http://192.168.15.101:5000/api/1.1/write/uo/uo1/100/16 
```

write to DOs or R1, R2 can write true/false/1/0

``` 
http://192.168.15.101:5000/api/1.1/write/do/r1/false/16
http://192.168.15.101:5000/api/1.1/write/do/r1/1/16
```

read UIS

``` 
http://192.168.15.101:5000/api/1.1/read/all/ui
```

read DIs

``` 
http://192.168.15.101:5000/api/1.1/read/all/di
```

### lora connect reset

// write the lora connect off

``` 
localhost:5000/api/1.1/write/do/lc/0/16
```

### RAW to Resistance Calculation (thermistor reading based on RAW 0-1)
```
T= (8.65943*((RAW/0.5555)/0.1774))/(9.89649 - (RAW/0.5555/0.1774))*1000
```