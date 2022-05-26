package api

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/bools"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"gorm.io/gorm"
	"math/bits"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Args struct {
	Sort                 string
	Order                string
	Offset               string
	Limit                string
	Search               string
	AskRefresh           string
	AskResponse          string
	Write                string
	ThingType            string
	UUID                 *string
	FlowNetworkUUID      string
	FlowNetworkCloneUUID *string
	WriteHistory         string
	WriteConsumer        string
	Field                string
	Value                string
	UpdateProducer       string
	CompactPayload       string
	CompactWithName      string
	FlowUUID             string
	Name                 *string
	StreamUUID           *string
	ProducerUUID         *string
	ConsumerUUID         *string
	WriterUUID           string
	AddToParent          string
	GlobalUUID           *string
	ClientId             *string
	SiteId               *string
	DeviceId             *string
	WriterThingClass     *string
	SourceUUID           *string
	ProducerThingUUID    *string
	WriterThingUUID      *string
	PluginName           bool
	TimestampGt          *string
	TimestampLt          *string
	WithFlowNetworks     bool
	WithStreams          bool
	WithStreamClones     bool
	WithProducers        bool
	WithConsumers        bool
	WithCommandGroups    bool
	WithWriters          bool
	WithWriterClones     bool
	Networks             bool
	WithDevices          bool
	WithPoints           bool
	WithPriority         bool
	WithTags             bool
	NetworkName          *string
	DeviceName           *string
	PointName            *string
	AddressUUID          *string
	AddressID            *string
	ObjectType           *string
	IoNumber             *string
	IdGt                 *string
	IsRemote             *bool
}

var ArgsType = struct {
	Sort                 string
	Order                string
	Offset               string
	Limit                string
	Search               string
	AskRefresh           string
	AskResponse          string
	Write                string
	ThingType            string
	UUID                 string
	FlowNetworkUUID      string
	FlowNetworkCloneUUID string
	WriteHistory         string
	WriteConsumer        string
	Field                string
	Name                 string
	Value                string
	UpdateProducer       string
	CompactPayload       string //for a point would be presentValue
	CompactWithName      string //for a point would be presentValue and pointName
	FlowUUID             string
	StreamUUID           string
	ProducerUUID         string
	ConsumerUUID         string
	WriterUUID           string
	AddToParent          string
	GlobalUUID           string
	ClientId             string
	SiteId               string
	DeviceId             string
	WriterThingClass     string
	SourceUUID           string
	ProducerThingUUID    string
	WriterThingUUID      string
	PluginName           string
	TimestampGt          string
	TimestampLt          string
	WithFlowNetworks     string
	WithStreams          string
	WithStreamClones     string
	WithProducers        string
	WithConsumers        string
	WithCommandGroups    string
	WithWriters          string
	WithWriterClones     string
	WithNetworks         string
	WithDevices          string
	WithPoints           string
	WithPriority         string
	WithTags             string
	NetworkName          string
	DeviceName           string
	PointName            string
	AddressUUID          string
	AddressID            string
	ObjectType           string
	IoNumber             string
	IdGt                 string
	IsRemote             string
}{
	Sort:                 "sort",
	Order:                "order",
	Offset:               "offset",
	Limit:                "limit",
	Search:               "search",
	AskRefresh:           "ask_refresh",
	AskResponse:          "ask_response",
	Write:                "write",      //consumer to write a value
	ThingType:            "thing_type", //the type of thing like a point
	UUID:                 "uuid",
	FlowNetworkUUID:      "flow_network_uuid",
	FlowNetworkCloneUUID: "flow_network_clone_uuid",
	WriteHistory:         "write_history",
	WriteConsumer:        "write_consumer",
	Field:                "field",
	Name:                 "name",
	Value:                "value",
	UpdateProducer:       "update_producer",
	CompactPayload:       "compact_payload",
	CompactWithName:      "compact_with_name",
	FlowUUID:             "flow_uuid",
	StreamUUID:           "stream_uuid",
	ProducerUUID:         "producer_uuid",
	ConsumerUUID:         "consumer_uuid",
	WriterUUID:           "writer_uuid",
	AddToParent:          "add_to_parent",
	GlobalUUID:           "global_uuid",
	ClientId:             "client_id",
	SiteId:               "site_id",
	DeviceId:             "device_id",
	WriterThingClass:     "writer_thing_class",
	SourceUUID:           "source_uuid",
	ProducerThingUUID:    "producer_thing_uuid",
	WriterThingUUID:      "writer_thing_uuid",
	PluginName:           "by_plugin_name",
	TimestampGt:          "timestamp_gt",
	TimestampLt:          "timestamp_lt",
	WithFlowNetworks:     "with_flow_networks",
	WithStreams:          "with_streams",
	WithStreamClones:     "with_stream_clones",
	WithProducers:        "with_producers",
	WithConsumers:        "with_consumers",
	WithCommandGroups:    "with_command_groups",
	WithWriters:          "with_writers",
	WithWriterClones:     "with_writer_clones",
	WithNetworks:         "with_networks",
	WithDevices:          "with_devices",
	WithPoints:           "with_points",
	WithPriority:         "with_priority",
	WithTags:             "with_tags",
	NetworkName:          "network_name",
	DeviceName:           "device_name",
	PointName:            "point_name",
	AddressUUID:          "address_uuid",
	AddressID:            "address_id",
	ObjectType:           "object_type",
	IoNumber:             "io_number",
	IdGt:                 "id_gt",
	IsRemote:             "is_remote",
}

var ArgsDefault = struct {
	Sort              string
	Order             string
	Offset            string
	Limit             string
	Search            string
	AskRefresh        string
	AskResponse       string
	Write             string
	ThingType         string
	FlowNetworkUUID   string
	Field             string
	Value             string
	UpdateProducer    string
	CompactPayload    string
	CompactWithName   string
	FlowUUID          string
	StreamUUID        string
	ProducerUUID      string
	ConsumerUUID      string
	WriterUUID        string
	AddToParent       string
	PluginName        string
	WithFlowNetworks  string
	WithStreams       string
	WithStreamClones  string
	WithProducers     string
	WithConsumers     string
	WithCommandGroups string
	WithWriters       string
	WithWriterClones  string
	WithNetworks      string
	WithDevices       string
	WithPoints        string
	WithPriority      string
	WithTags          string
	NetworkName       string
	DeviceName        string
	PointName         string
}{
	Sort:              "ID",
	Order:             "DESC", //ASC or DESC
	Offset:            "0",
	Limit:             "25",
	Search:            "",
	AskRefresh:        "false",
	AskResponse:       "false",
	Write:             "false",
	ThingType:         "point",
	FlowNetworkUUID:   "",
	Field:             "name",
	Value:             "",
	UpdateProducer:    "false",
	CompactPayload:    "false",
	CompactWithName:   "false",
	FlowUUID:          "",
	ProducerUUID:      "",
	ConsumerUUID:      "",
	WriterUUID:        "",
	AddToParent:       "",
	PluginName:        "false",
	WithFlowNetworks:  "false",
	WithStreams:       "false",
	WithStreamClones:  "false",
	WithProducers:     "false",
	WithConsumers:     "false",
	WithCommandGroups: "false",
	WithWriters:       "false",
	WithWriterClones:  "false",
	WithNetworks:      "false",
	WithDevices:       "false",
	WithPoints:        "false",
	WithPriority:      "false",
	WithTags:          "false",
	NetworkName:       "",
	DeviceName:        "",
	PointName:         "",
}

func responseHandler(body interface{}, err error, ctx *gin.Context) {
	// TODO: add other custom errors
	if err == nil {
		ctx.JSON(http.StatusOK, body)
	} else {
		switch err {
		case gorm.ErrRecordNotFound:
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 404, err.Error())
			ctx.JSON(http.StatusNotFound, interfaces.Message{Message: message})
		case gorm.ErrInvalidTransaction,
			gorm.ErrNotImplemented,
			gorm.ErrMissingWhereClause,
			gorm.ErrUnsupportedRelation,
			gorm.ErrPrimaryKeyRequired,
			gorm.ErrModelValueRequired,
			gorm.ErrInvalidData,
			gorm.ErrUnsupportedDriver,
			gorm.ErrRegistered,
			gorm.ErrInvalidField,
			gorm.ErrEmptySlice,
			gorm.ErrDryRunModeUnsupported,
			gorm.ErrInvalidDB,
			gorm.ErrInvalidValue,
			gorm.ErrInvalidValueOfLength:
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 500, err.Error())
			ctx.JSON(http.StatusInternalServerError, interfaces.Message{Message: message})
		default:
			message := fmt.Sprintf("%s %s [%d]: %s", ctx.Request.Method, ctx.Request.URL, 400, err.Error())
			ctx.JSON(http.StatusBadRequest, interfaces.Message{Message: message})
		}
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

func getBODYScheduleData(ctx *gin.Context) (dto *model.ScheduleData, err error) {
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

func getMapBody(ctx *gin.Context) (dto *map[string]interface{}, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYNetwork(ctx *gin.Context) (dto *model.Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyHistory(ctx *gin.Context) (dto *model.ProducerHistory, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyBulkHistory(ctx *gin.Context) (dto []*model.ProducerHistory, err error) {
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

func getBodySyncWriter(ctx *gin.Context) (dto *model.SyncWriter, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySyncCOV(ctx *gin.Context) (dto *model.SyncCOV, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodySyncWriterAction(ctx *gin.Context) (dto *model.SyncWriterAction, err error) {
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

func getBODYWriterBulk(ctx *gin.Context) (dto []*model.WriterBulkBody, err error) {
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

func getBodySyncStream(ctx *gin.Context) (dto *model.SyncStream, err error) {
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

func getBODYPointMapping(ctx *gin.Context) (dto *model.PointMapping, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYBulkPoints(ctx *gin.Context) (dto []*model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPoint(ctx *gin.Context) (dto *model.Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYPointWriter(ctx *gin.Context) (dto *model.PointWriter, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyTag(ctx *gin.Context) (dto *model.Tag, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getP2PBody(ctx *gin.Context) (dto *model.P2PBody, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveGlobalUUID(ctx *gin.Context) string {
	return ctx.Param("global_uuid")
}

func resolveID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

func resolveWriterUUID(ctx *gin.Context) string {
	return ctx.Param("writer_uuid")
}

func resolveSourceUUID(ctx *gin.Context) string {
	return ctx.Param("source_uuid")
}

func resolveWriterThingClass(ctx *gin.Context) string {
	return ctx.Param("writer_thing_class")
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

func resolveName(ctx *gin.Context) string {
	return ctx.Param("name")
}

func resolvePluginName(ctx *gin.Context) string {
	return ctx.Param("plugin_name")
}

func getRubixPingDevice(ctx *gin.Context) (dto *Ping, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func toBool(value string) (bool, error) {
	if value == "" {
		return false, nil
	} else {
		c, err := bools.Boolean(value)
		return c, err
	}
}
