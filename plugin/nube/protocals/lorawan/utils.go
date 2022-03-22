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
	o := mArr.Get(0)
	a := mArr.Get(1)
	return o.(string), a.(string)
}
