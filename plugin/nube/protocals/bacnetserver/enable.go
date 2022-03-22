package main

import (
	"github.com/NubeIO/flow-framework/api"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nrest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube_api"
	nube_api_bacnetserver "github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube_api/bacnetserver"
	"github.com/labstack/gommon/log"
	"time"
)

//inc generic reset client
var reqType = &nrest.ReqType{
	BaseUri: nube_api.BaseURL,
	Service: "bacnet-server",
	LogPath: "helpers.nrest.bacnet.server",
	Port:    nube_api.DefaultPortBacnet,
}

//api options
var options = &nrest.ReqOpt{
	Timeout:          500 * time.Second,
	RetryCount:       0,
	RetryWaitTime:    0 * time.Second,
	RetryMaxWaitTime: 0,
	//Headers:          map[string]interface{}{"Authorization": nubeApi.RubixToken},
}

//inc nube rest client
var nubeApi = &nube_api.NubeRest{
	Rest:          reqType,
	RubixPort:     nube_api.DefaultRubixService,
	RubixUsername: "",
	RubixPassword: "",
	UseRubixProxy: false,
}

var nubeClient = nube_api.New(nubeApi)

var bacnetClient = &nube_api_bacnetserver.RestClient{
	NubeRest: nubeClient,
	Options:  options,
}

// Enable implements plugin.Plugin
func (inst *Instance) Enable() error {
	inst.enabled = true
	inst.setUUID()
	inst.BusServ()
	q, err := inst.db.GetNetworkByPlugin(inst.pluginUUID, api.Args{})
	if q != nil {
		inst.networkUUID = q.UUID
	} else {
		inst.networkUUID = "NA"
	}
	if err != nil {
		log.Error("error on enable bacnetserver-plugin")
	}
	return nil
}

// Disable implements plugin.Disable
func (inst *Instance) Disable() error {
	inst.enabled = false
	return nil
}
