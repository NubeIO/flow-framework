package database

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nuuid"
	"github.com/NubeIO/flow-framework/utils/structs"
	"gorm.io/gorm"

	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
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
		uuid := nuuid.MakeTopicUUID("")
		return fmt.Sprintf("n_%s", truncateString(uuid, 8))
	}
	return name
}

func checkTransport(t string) (string, error) {
	if t == "" {
		return model.TransType.IP, nil
	}
	i := structs.ArrayValues(model.TransType)
	if !structs.ArrayContains(i, t) {
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

func checkHistoryType(t string) (model.HistoryType, error) {
	if t == "" {
		return model.HistoryTypeInterval, nil
	}
	historyType := model.HistoryType(t)
	if _, ok := model.HistoryTypeMap[historyType]; !ok {
		return "", errors.New("please provide a valid history type ie: COV , INTERVAL or COV_AND_INTERVAL")
	}
	return historyType, nil
}

func checkHistoryCovType(t string) bool {
	if t == "" {
		return false
	}
	historyType := model.HistoryType(t)
	if _, ok := model.HistoryTypeCovMap[historyType]; !ok {
		return false
	}
	return true
}

func (d *GormDatabase) deleteResponseBuilder(query *gorm.DB) (bool, error) {
	if query.Error != nil {
		return false, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return false, gorm.ErrRecordNotFound
	} else {
		return true, nil
	}
}
