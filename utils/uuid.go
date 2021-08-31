package utils
import (
	"fmt"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
)


func MakeUUID() (string, error) {
	return uuid.MakeUUID()
}


func MakeTopicUUID(attribute string) string {
	u, _ := uuid.MakeUUID()
	divider := "_"
	fln := "fln" //flow network
	rfn := "rfn" //flow network
	plg := "plg" //plugin
	net := "net" //network
	dev := "dev" //device
	pnt := "pnt" //point
	job := "job" //job
	str := "str" //stream gateway
	stl := "stl" //list of flow network gateway
	pro := "pro" //producers
	prl := "prl" //producersList
	prh := "prh" //producer history
	sus := "sub" //consumers
	sul := "sul" //consumersList
	alt := "alt" //alerts
	cmd := "cmd" //command
	rub := "rbx" //rubix uuid
	rxg := "rxg" //rubix global uuid

	switch attribute {
	case model.CommonNaming.Plugin:
		return fmt.Sprintf("%s%s%s", plg, divider, u)
	case model.CommonNaming.FlowNetwork:
		return fmt.Sprintf("%s%s%s", fln, divider, u)
	case model.CommonNaming.RemoteFlowNetwork:
		return fmt.Sprintf("%s%s%s", rfn, divider, u)
	case model.CommonNaming.Network:
		return fmt.Sprintf("%s%s%s", net, divider, u)
	case model.CommonNaming.Device:
		return fmt.Sprintf("%s%s%s", dev, divider, u)
	case model.CommonNaming.Point:
		return fmt.Sprintf("%s%s%s", pnt, divider, u)
	case model.CommonNaming.Stream:
		return fmt.Sprintf("%s%s%s", str, divider, u)
	case model.CommonNaming.StreamList:
		return fmt.Sprintf("%s%s%s", stl, divider, u)
	case model.CommonNaming.Job:
		return fmt.Sprintf("%s%s%s", job, divider, u)
	case model.CommonNaming.Producer:
		return fmt.Sprintf("%s%s%s", pro, divider, u)
	case model.CommonNaming.WriterCopy:
		return fmt.Sprintf("%s%s%s", prl, divider, u)
	case model.CommonNaming.ProducerHistory:
		return fmt.Sprintf("%s%s%s", prh, divider, u)
	case model.CommonNaming.Consumer:
		return fmt.Sprintf("%s%s%s", sus, divider, u)
	case model.CommonNaming.Writer:
		return fmt.Sprintf("%s%s%s", sul, divider, u)
	case model.CommonNaming.Alert:
		return fmt.Sprintf("%s%s%s", alt, divider, u)
	case model.CommonNaming.CommandGroup:
		return fmt.Sprintf("%s%s%s", cmd, divider, u)
	case model.CommonNaming.Rubix:
		return fmt.Sprintf("%s%s%s", rub, divider, u)
	case model.CommonNaming.RubixGlobal:
		return fmt.Sprintf("%s%s%s", rxg, divider, u)

	}
	return u
}

