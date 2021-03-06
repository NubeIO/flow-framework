package helpers

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils/nuuid"
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

func NameIsNil() string {
	uuid := nuuid.MakeTopicUUID("")
	return fmt.Sprintf("n_%s", truncateString(uuid, 8))
}
