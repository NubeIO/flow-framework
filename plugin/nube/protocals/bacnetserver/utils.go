package main

import (
	"github.com/NubeIO/flow-framework/utils"
	"strings"
)

func getPointAddr(s string) (objType, addr string) {
	mArr := utils.NewArray()
	ss := strings.Split(s, "-")
	for _, e := range ss {
		if e != "" {
			mArr.Add(e)
		}
	}
	const objectType = 0
	const address = 1
	o := mArr.Get(objectType)
	a := mArr.Get(address)
	return o.(string), a.(string)
}
