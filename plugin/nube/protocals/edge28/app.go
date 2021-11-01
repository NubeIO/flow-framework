package main

var Rls = []string{"R1", "R2"}
var DOs = []string{"DO1", "DO2", "DO3", "DO4", "DO5"}
var UOs = []string{"UO1", "UO2", "UO3", "UO4", "UO5", "UO6", "UO7"}
var UIs = []string{"UI1", "UI2", "UI3", "UI4", "UI5", "UI6", "UI7"}
var DIs = []string{"DI1", "DI2", "DI3", "DI4", "DI5", "DI6", "DI7"}

var pointList = struct {
	R1  string `json:"R1"`
	R2  string `json:"R2"`
	DO1 string `json:"DO1"`
	DO2 string `json:"DO2"`
	DO3 string `json:"DO3"`
	DO4 string `json:"DO4"`
	DO5 string `json:"DO5"`
	UO1 string `json:"UO1"`
	UO2 string `json:"UO2"`
	UO3 string `json:"UO3"`
	UO4 string `json:"UO4"`
	UO5 string `json:"UO5"`
	UO6 string `json:"UO6"`
	UO7 string `json:"UO7"`
	UI1 string `json:"UI1"`
	UI2 string `json:"UI2"`
	UI3 string `json:"UI3"`
	UI4 string `json:"UI4"`
	UI5 string `json:"UI5"`
	UI6 string `json:"UI6"`
	UI7 string `json:"UI7"`
	DI1 string `json:"DI1"`
	DI2 string `json:"DI2"`
	DI3 string `json:"DI3"`
	DI4 string `json:"DI4"`
	DI5 string `json:"DI5"`
	DI6 string `json:"DI6"`
	DI7 string `json:"DI7"`
}{
	R1:  "R1",
	R2:  "R2",
	DO1: "DO1",
	DO2: "DO2",
	DO3: "DO3",
	DO4: "DO4",
	DO5: "DO5",
	UO1: "UO1",
	UO2: "UO2",
	UO3: "UO3",
	UO4: "UO4",
	UO5: "UO5",
	UO6: "UO6",
	UO7: "UO7",
	UI1: "UI1",
	UI2: "UI2",
	UI3: "UI3",
	UI4: "UI4",
	UI5: "UI5",
	UI6: "UI6",
	UI7: "UI7",
	DI1: "DI1",
	DI2: "DI2",
	DI3: "DI3",
	DI4: "DI4",
	DI5: "DI5",
	DI6: "DI6",
	DI7: "DI7",
}

func pointsAll() []string {
	out := append(Rls, DOs...)
	out = append(out, UOs...)
	out = append(out, DIs...)
	out = append(out, UIs...)
	return out
}

var UOTypes = struct {
	//RAW  string
	DIGITAL string
	PERCENT string
	VOLTSDC string
	//MILLIAMPS  string
}{
	//RAW:  "RAW",
	DIGITAL: "DIGITAL",
	PERCENT: "PERCENT",
	VOLTSDC: "0-10VDC",
	//MILLIAMPS:  "4-20mA",
}

var UITypes = struct {
	RAW              string
	DIGITAL          string
	PERCENT          string
	VOLTSDC          string
	MILLIAMPS        string
	RESISTANCE       string
	THERMISTOR10KT2  string
	THERMISTOR10KT3  string
	THERMISTOR20KT1  string
	THERMISTORPT100  string
	THERMISTORPT1000 string
}{
	RAW:              "RAW",
	DIGITAL:          "DIGITAL",
	PERCENT:          "PERCENT",
	VOLTSDC:          "0-10VDC",
	MILLIAMPS:        "4-20mA",
	RESISTANCE:       "RESISTANCE",
	THERMISTOR10KT2:  "THERMISTOR_10K_TYPE2",
	THERMISTOR10KT3:  "THERMISTOR_10K_TYPE3",
	THERMISTOR20KT1:  "THERMISTOR_20K_TYPE1",
	THERMISTORPT100:  "THERMISTOR_PT100",
	THERMISTORPT1000: "THERMISTOR_PT1000",
}
