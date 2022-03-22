package main

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type Network struct {
	NetworkUUID           string `json:"network_uuid"`
	NetworkName           string `json:"network_name"`
	NetworkIp             string `json:"network_ip"`
	NetworkMask           int    `json:"network_mask"`
	NetworkPort           int    `json:"network_port"`
	NetworkDeviceObjectId int    `json:"network_device_object_id"`
	NetworkDeviceName     string `json:"network_device_name"`
	InterfaceName         string `json:"interface_name"`
}

type Networks struct {
	Networks []struct {
		NetworkUuid string `json:"network_uuid"`
	} `json:"networks"`
}

type Device struct {
	DeviceUUID     string `json:"device_uuid"`
	DeviceName     string `json:"device_name"`
	DeviceIp       string `json:"device_ip"`
	DeviceEnable   bool   `json:"device_enable"`
	DeviceMask     int    `json:"device_mask"`
	DevicePort     int    `json:"device_port"`
	DeviceMac      int    `json:"device_mac"`
	DeviceObjectId int    `json:"device_object_id"`
	NetworkNumber  int    `json:"network_number"`
	TypeMstp       bool   `json:"type_mstp"`
	SupportsRpm    bool   `json:"supports_rpm"`
	SupportsWpm    bool   `json:"supports_wpm"`
	NetworkUUID    string `json:"network_uuid"`
}

type Point struct {
	PointUUID       string `json:"point_uuid"`
	PointName       string `json:"point_name"`
	PointEnable     bool   `json:"point_enable"`
	PointObjectId   int    `json:"point_object_id"`
	PointObjectType string `json:"point_object_type"`
	DeviceUUID      string `json:"device_uuid"`
}

type Whois struct {
	NetworkNumber         int  `json:"network_number"`
	IsMstp                bool `json:"is_mstp"`
	Whois                 bool `json:"whois"`
	GlobalBroadcast       bool `json:"global_broadcast"`
	FullRange             bool `json:"full_range"`
	RangeStart            int  `json:"range_start"`
	RangeEnd              int  `json:"range_end"`
	ShowSupportedServices bool `json:"show_supported_services"`
	AddDevices            bool `json:"add_devices"`
}

type ReleaseOverride struct {
	Priority int  `json:"priority"`
	Feedback bool `json:"feedback"`
	Timeout  int  `json:"timeout"`
}

type PointWrite struct {
	PointWriteValue int  `json:"point_write_value"`
	Priority        int  `json:"priority"`
	Feedback        bool `json:"feedback"`
	Timeout         int  `json:"timeout"`
}

type PointRead struct {
	GetPriority bool `json:"get_priority"`
	Timeout     int  `json:"timeout"`
}

type PointsDiscover struct {
	AddPoints bool `json:"add_points"`
	Discovery bool `json:"discovery"`
	Timeout   int  `json:"timeout"`
}

const (
	schemaNetwork = "/schema/network"
	schemaDevice  = "/schema/device"
	schemaPoint   = "/schema/point"

	wizard         = "/wizard"
	ping           = "/ping"
	pingBacnet     = "/api/system/ping"
	point          = "/point"
	pointBacnet    = "/api/bm/point"
	points         = "/points"
	pointsBacnet   = "/api/bm/points"
	network        = "/network"
	networkBacnet  = "/api/bm/network"
	networks       = "/networks"
	networksBacnet = "/api/bm/networks/true"
	device         = "/device"
	devices        = "/devices"
	deviceBacnet   = "/api/bm/device"
	devicesBacnet  = "/api/bm/devices"
)

func getBODYAddNetwork(ctx *gin.Context) (dto *Network, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYAddDevice(ctx *gin.Context) (dto *Device, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func getBODYAddPoint(ctx *gin.Context) (dto *Point, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func resolveObject(ctx *gin.Context) string {
	return ctx.Param("object")
}

func resolveAddress(ctx *gin.Context) string {
	return ctx.Param("address")
}

func resolveUUID(ctx *gin.Context) string {
	return ctx.Param("uuid")
}

// RegisterWebhook implements plugin.Webhooker
func (i *Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	i.basePath = basePath
	mux.POST(wizard, func(ctx *gin.Context) {
		body, err := getBODYAddNetwork(ctx)
		_, err = i.wizard(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, "ok")
			return
		}
	})
	mux.GET(ping, func(ctx *gin.Context) {
		rt.Method = nrest.GET
		rt.Path = pingBacnet
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		if err != nil {
			ctx.JSON(code, err)
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.GET(devices, func(ctx *gin.Context) {
		rt.Method = nrest.GET
		rt.Path = devicesBacnet + "/true"
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		if err != nil {
			ctx.JSON(code, err)
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.GET(points, func(ctx *gin.Context) {
		rt.Method = nrest.GET
		rt.Path = pointsBacnet
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		if err != nil {
			ctx.JSON(code, err)
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.GET(networks, func(ctx *gin.Context) {
		rt.Method = nrest.GET
		rt.Path = networksBacnet
		res, code, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		if err != nil {
			ctx.JSON(code, err)
			return
		} else {
			ctx.JSON(code, res.AsJsonNoErr())
			return
		}
	})
	mux.POST(network, func(ctx *gin.Context) {
		body, _ := getBODYAddNetwork(ctx)
		net, res, code, err := i.addNetwork(body)
		if net != nil {
			ctx.JSON(code, err)
			return
		}
		if res != nil {
			ctx.JSON(code, res)
			return
		}
		if err != nil {
			ctx.JSON(code, err)
			return
		}
	})
	mux.POST(device, func(ctx *gin.Context) {
		body, _ := getBODYAddDevice(ctx)
		add, err := i.addDevice(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.POST(point, func(ctx *gin.Context) {
		body, _ := getBODYAddPoint(ctx)
		add, err := i.addPoint(body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.PATCH("network/:uuid", func(ctx *gin.Context) {
		body, _ := getBODYAddPoint(ctx)
		uuid := resolveUUID(ctx)
		add, err := i.patchPoint(body, uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.PATCH("device/:uuid", func(ctx *gin.Context) {
		body, _ := getBODYAddDevice(ctx)
		uuid := resolveUUID(ctx)
		add, err := i.patchDevice(body, uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.PATCH("point/:uuid", func(ctx *gin.Context) {
		body, _ := getBODYAddPoint(ctx)
		uuid := resolveUUID(ctx)
		add, err := i.patchPoint(body, uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.DELETE("network/:uuid", func(ctx *gin.Context) {
		uuid := resolveUUID(ctx)
		add, err := i.deleteNetwork(uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.DELETE("device/:uuid", func(ctx *gin.Context) {
		uuid := resolveUUID(ctx)
		add, err := i.deleteDevice(uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.DELETE("point/:uuid", func(ctx *gin.Context) {
		uuid := resolveUUID(ctx)
		add, err := i.deletePoint(uuid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		} else {
			ctx.JSON(http.StatusOK, add)
			return
		}
	})
	mux.DELETE(networks, func(ctx *gin.Context) {
		//getNetworks
		rt.Method = nrest.GET
		rt.Path = networksBacnet
		res, _, err := nrest.DoHTTPReq(rt, &nrest.ReqOpt{})
		nets := new(Networks)
		err = res.ToInterface(nets)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errors.New("failed to find get networks"))
		}
		count := 0
		for _, a := range nets.Networks {
			rt.Method = nrest.DELETE
			rt.Path = networkBacnet + "/" + a.NetworkUuid
			fmt.Println(rt.Path)
			count++
			nrest.DoHTTPReq(rt, &nrest.ReqOpt{})

		}
		r := fmt.Sprintf("deleted %d number of networks", count)
		ctx.JSON(http.StatusOK, r)
	})
}
