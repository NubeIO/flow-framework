package utils
import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
)


func MakeUUID() (string, error) {
	return uuid.MakeUUID()
}


func MakeTopicUUID(attribute string) (string, error) {
	u, err := uuid.MakeUUID()
	divider := "_"
	net := "net"
	dev := "dev"
	pnt := "pnt"
	job := "pnt"
	gtw := "gtw"
	sub := "sub"
	rip := "rip"
	alm := "alm"
	switch attribute {
	case model.CommonNaming.Network:
		return fmt.Sprintf("%s%s%s", net, divider, u), err
	case model.CommonNaming.Device:
		fmt.Println("ssssssss")
		return fmt.Sprintf("%s%s%s", dev, divider, u), err
	case model.CommonNaming.Point:
		return fmt.Sprintf("%s%s%s", pnt, divider, u), err
	case model.CommonNaming.Gateway:
		return fmt.Sprintf("%s%s%s", gtw, divider, u), err
	case model.CommonNaming.Job:
		return fmt.Sprintf("%s%s%s", job, divider, u), err
	case model.CommonNaming.Subscriber:
		return fmt.Sprintf("%s%s%s", sub, divider, u), err
	case model.CommonNaming.Subscription:
		return fmt.Sprintf("%s%s%s", rip, divider, u), err
	case model.CommonNaming.Alarm:
		return fmt.Sprintf("%s%s%s", alm, divider, u), err

	}
	fmt.Println("here")
	return uuid.MakeUUID()
}

