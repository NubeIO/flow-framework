package api

import (
	"errors"
	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/bools"
	"math/bits"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Args struct {
	Sort             string
	Order            string
	Offset           string
	Limit            string
	Search           string
	WithChildren     string
	WithPoints       string
	WithParent       string
	AskRefresh       string
	AskResponse      string
	Write            string
	ThingType        string
	FlowNetworkUUID  string
	WriteHistory     string
	WriteConsumer    string
	Field            string
	Value            string
	UpdateProducer   string
	CompactPayload   string
	CompactWithName  string
	FlowUUID         string
	StreamUUID       string
	ProducerUUID     string
	ConsumerUUID     string
	WriterUUID       string
	AddToParent      string
	FlowNetworks     bool
	Streams          bool
	StreamClones     bool
	Producers        bool
	Consumers        bool
	CommandGroups    bool
	Writers          bool
	WriterClones     bool
	Networks         bool
	Devices          bool
	Points           bool
	Priority         bool
	IpConnection     bool
	SerialConnection bool
	Tags             bool
	GlobalUUID       *string
	ClientId         *string
	SiteId           *string
	DeviceId         *string
	PluginName       bool
	TimestampGt      *string
	TimestampLt      *string
}

var ArgsType = struct {
	Sort             string
	Order            string
	Offset           string
	Limit            string
	Search           string
	WithChildren     string
	WithPoints       string
	WithParent       string
	AskRefresh       string
	AskResponse      string
	Write            string
	ThingType        string
	FlowNetworkUUID  string
	WriteHistory     string
	WriteConsumer    string
	Field            string
	Value            string
	UpdateProducer   string
	CompactPayload   string //for a point would be presentValue
	CompactWithName  string //for a point would be presentValue and pointName
	FlowUUID         string
	StreamUUID       string
	ProducerUUID     string
	ConsumerUUID     string
	WriterUUID       string
	AddToParent      string
	FlowNetworks     string
	Streams          string
	StreamClones     string
	Producers        string
	Consumers        string
	CommandGroups    string
	Writers          string
	WriterClones     string
	Networks         string
	Devices          string
	Points           string
	Priority         string
	IpConnection     string
	SerialConnection string
	Tags             string
	GlobalUUID       string
	ClientId         string
	SiteId           string
	DeviceId         string
	PluginName       string
	TimestampGt      string
	TimestampLt      string
}{
	Sort:             "sort",
	Order:            "order",
	Offset:           "offset",
	Limit:            "limit",
	Search:           "search",
	WithChildren:     "with_children",
	WithPoints:       "with_points",
	WithParent:       "with_parent",
	AskRefresh:       "ask_refresh",
	AskResponse:      "ask_response",
	Write:            "write",             //consumer to write a value
	ThingType:        "thing_type",        //the type of thing like a point
	FlowNetworkUUID:  "flow_network_uuid", //the type of thing like a point
	WriteHistory:     "write_history",
	WriteConsumer:    "write_consumer",
	Field:            "field",
	Value:            "value",
	UpdateProducer:   "update_producer",
	CompactPayload:   "compact_payload",
	CompactWithName:  "compact_with_name",
	FlowUUID:         "flow_uuid",
	StreamUUID:       "stream_uuid",
	ProducerUUID:     "producer_uuid",
	ConsumerUUID:     "consumer_uuid",
	WriterUUID:       "writer_uuid",
	AddToParent:      "add_to_parent",
	FlowNetworks:     "flow_networks",
	Streams:          "streams",
	StreamClones:     "stream_clones",
	Producers:        "producers",
	Consumers:        "consumers",
	CommandGroups:    "command_groups",
	Writers:          "writers",
	WriterClones:     "writer_clones",
	Networks:         "networks",
	Devices:          "devices",
	Points:           "points",
	Priority:         "priority",
	IpConnection:     "ip_connection",
	SerialConnection: "serial_connection",
	Tags:             "tags",
	GlobalUUID:       "global_uuid",
	ClientId:         "client_id",
	SiteId:           "site_id",
	DeviceId:         "device_id",
	PluginName:       "by_plugin_name",
	TimestampGt:      "timestamp_gt",
	TimestampLt:      "timestamp_lt",
}

var ArgsDefault = struct {
	Sort             string
	Order            string
	Offset           string
	Limit            string
	Search           string
	WithChildren     string
	WithPoints       string
	WithParent       string
	AskRefresh       string
	AskResponse      string
	Write            string
	ThingType        string
	FlowNetworkUUID  string
	Field            string
	Value            string
	UpdateProducer   string
	CompactPayload   string
	CompactWithName  string
	FlowUUID         string
	StreamUUID       string
	ProducerUUID     string
	ConsumerUUID     string
	WriterUUID       string
	AddToParent      string
	FlowNetworks     string
	Streams          string
	StreamClones     string
	Producers        string
	Consumers        string
	CommandGroups    string
	Writers          string
	WriterClones     string
	Networks         string
	Devices          string
	Points           string
	Priority         string
	IpConnection     string
	SerialConnection string
	Tags             string
	PluginName       string
}{
	Sort:             "ID",
	Order:            "DESC", //ASC or DESC
	Offset:           "0",
	Limit:            "25",
	Search:           "",
	WithChildren:     "false",
	WithPoints:       "false",
	WithParent:       "false",
	AskRefresh:       "false",
	AskResponse:      "false",
	Write:            "false",
	ThingType:        "point",
	FlowNetworkUUID:  "",
	Field:            "name",
	Value:            "",
	UpdateProducer:   "false",
	CompactPayload:   "false",
	CompactWithName:  "false",
	FlowUUID:         "",
	StreamUUID:       "",
	ProducerUUID:     "",
	ConsumerUUID:     "",
	WriterUUID:       "",
	AddToParent:      "",
	FlowNetworks:     "false",
	Streams:          "false",
	StreamClones:     "false",
	Producers:        "false",
	Consumers:        "false",
	CommandGroups:    "false",
	Writers:          "false",
	WriterClones:     "false",
	Networks:         "false",
	Devices:          "false",
	Points:           "false",
	Priority:         "false",
	IpConnection:     "false",
	SerialConnection: "false",
	Tags:             "false",
	PluginName:       "false",
}

func reposeHandler(body interface{}, err error, ctx *gin.Context) {
	if err != nil {
		if body == nil {
			ctx.JSON(404, "unknown error")
		} else {
			ctx.JSON(404, err.Error())
		}
	} else {
		ctx.JSON(200, body)
	}
}

func withID(ctx *gin.Context, name string, f func(id uint)) {
	if id, err := strconv.ParseUint(ctx.Param(name), 10, bits.UintSize); err == nil {
		f(uint(id))
	} else {
		ctx.AbortWithError(400, errors.New("invalid id"))
	}
}

func getBODYSchedule(ctx *gin.Context) (dto *model.Schedule, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyLocalStorageFlowNetwork(ctx *gin.Context) (dto *model.LocalStorageFlowNetwork, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYFlowNetwork(ctx *gin.Context) (dto *model.FlowNetwork, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYNetwork(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYHistory(ctx *gin.Context) (dto *model.ProducerHistory, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYBulkHistory(ctx *gin.Context) (dto []*model.ProducerHistory, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYDevice(ctx *gin.Context) (dto *model.Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYProducer(ctx *gin.Context) (dto *model.Producer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYConsumer(ctx *gin.Context) (dto *model.Consumer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriter(ctx *gin.Context) (dto *model.Writer, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPlugin(ctx *gin.Context) (dto *model.PluginConf, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYMqttConnection(ctx *gin.Context) (dto *model.MqttConnection, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYIntegration(ctx *gin.Context) (dto *model.Integration, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterBody(ctx *gin.Context) (dto *model.WriterBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterBulk(ctx *gin.Context) (dto []*model.WriterBulk, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYWriterClone(ctx *gin.Context) (dto *model.WriterClone, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYStream(ctx *gin.Context) (dto *model.Stream, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYCommandGroup(ctx *gin.Context) (dto *model.CommandGroup, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYJobs(ctx *gin.Context) (dto *model.Job, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPoint(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyTag(ctx *gin.Context) (dto *model.Tag, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

func resolveProducerUUID(ctx *gin.Context) string {
	return ctx.Param("producer_uuid")
}

func getTagParam(ctx *gin.Context) string {
	return ctx.Param("tag")
}

func resolvePath(ctx *gin.Context) string {
	return ctx.Param("path")
}

func resolvePluginUUID(ctx *gin.Context) string {
	return ctx.Param("plugin_uuid")
}

func resolvePluginName(ctx *gin.Context) string {
	return ctx.Param("plugin_name")
}

func toBool(value string) (bool, error) {
	if value == "" {
		return false, nil
	} else {
		c, err := bools.Boolean(value)
		return c, err
	}
}
