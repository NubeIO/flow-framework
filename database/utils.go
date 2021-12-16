package database

import (
	"errors"
	"fmt"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/utils"
)

func truncateString(str string, num int) string {
	ret := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		ret = str[0:num] + ""
	}
	return ret
}

func typeIsNil(t string, use string) string {
	if t == "" {
		return use
	}
	return t
}

func pluginIsNil(name string) string {
	if name == "" {
		return "system"
	}
	return name
}

func nameIsNil(name string) string {
	if name == "" {
		uuid := utils.MakeTopicUUID("")
		return fmt.Sprintf("n_%s", truncateString(uuid, 8))
	}
	return name
}

func checkTransport(t string) (string, error) {
	if t == "" {
		return model.TransType.IP, nil
	}
	i := utils.ArrayValues(model.TransType)
	if !utils.ArrayContains(i, t) {
		return "", errors.New("please provide a valid transport type ie: ip or serial")
	}
	return t, nil
}

func checkObjectType(t string) (model.ObjectType, error) {
	if t == "" {
		return model.ObjTypeAnalogValue, nil
	}
	objType := model.ObjectType(t)
	if _, ok := model.ObjectTypesMap[objType]; !ok {
		return "", errors.New("please provide a valid object type ie: analogInput or readCoil")
	}
	return objType, nil
}
