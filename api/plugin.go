package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/NubeDev/location"
	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin"
	"github.com/NubeIO/flow-framework/plugin/compat"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// The PluginDatabase interface for encapsulating database access.
type PluginDatabase interface {
	GetPlugins() ([]*model.PluginConf, error)
	GetPluginByPath(path string) (*model.PluginConf, error)
	GetPlugin(uuid string) (*model.PluginConf, error)
	UpdatePluginConf(p *model.PluginConf) error
}

// The PluginAPI provides handlers for managing plugins.
type PluginAPI struct {
	Notifier Notifier
	Manager  *plugin.Manager
	DB       PluginDatabase
}

func (c *PluginAPI) GetPlugin(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	q, err := c.DB.GetPlugin(uuid)
	responseHandler(q, err, ctx)
}

func (c *PluginAPI) GetPluginByPath(ctx *gin.Context) {
	path := resolvePath(ctx)
	q, err := c.DB.GetPluginByPath(path)
	responseHandler(q, err, ctx)
}

// GetPlugins returns all plugins a user has.
func (c *PluginAPI) GetPlugins(ctx *gin.Context) {
	plugins, err := c.DB.GetPlugins()
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	result := make([]model.PluginConfExternal, 0)
	for _, conf := range plugins {
		if inst, err := c.Manager.Instance(conf.UUID); err == nil {
			info := c.Manager.PluginInfo(conf.ModulePath)
			result = append(result, model.PluginConfExternal{
				UUID:         conf.UUID,
				Name:         info.String(),
				ModulePath:   conf.ModulePath,
				Author:       info.Author,
				Website:      info.Website,
				License:      info.License,
				Enabled:      conf.Enabled,
				HasNetwork:   conf.HasNetwork,
				Capabilities: inst.Supports().Strings(),
			})
		}
	}
	responseHandler(result, err, ctx)
}

//buildUUID a way to query a plugin by its name or uuid
func (c *PluginAPI) buildUUID(ctx *gin.Context) string {
	nameOrUUID := resolveID(ctx) //system?by_plugin_name=true is passed in then enable plugin by its name
	uuid := ""
	args := buildPluginArgs(ctx)
	if args.PluginName {
		path, err := c.DB.GetPluginByPath(nameOrUUID)
		if err != nil {
			responseHandler("err: no plugin with that name was found", nil, ctx)
		} else {
			uuid = path.UUID
		}
	} else {
		uuid = resolveID(ctx)
	}
	if uuid == "" {
		responseHandler("err: no valid uuid found", nil, ctx)
	}
	return uuid
}

// EnablePluginByUUID enables a plugin.
func (c *PluginAPI) EnablePluginByUUID(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	body, err := getBODYPlugin(ctx)
	if err != nil {
		responseHandler("error on body", err, ctx)
	}
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		responseHandler("unknown plugin", err, ctx)
		return
	}
	_, err = c.Manager.Instance(uuid)
	if err != nil {
		responseHandler("plugin not found", err, ctx)
		return
	}
	if err := c.Manager.SetPluginEnabled(uuid, body.Enabled); err == plugin.ErrAlreadyEnabledOrDisabled {
		responseHandler(nil, err, ctx)
		return
	} else if err != nil {
		responseHandler(nil, err, ctx)
		return
	}
	if body.Enabled {
		responseHandler(map[string]string{"state": "enabled"}, err, ctx)
	} else {
		responseHandler(map[string]string{"state": "disabled"}, err, ctx)
	}
}

// RestartPlugin enables a plugin.
func (c *PluginAPI) RestartPlugin(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		responseHandler("unknown plugin", err, ctx)
		return
	}
	_, err = c.Manager.Instance(uuid)
	if err != nil {
		responseHandler("plugin not found", err, ctx)
		return
	}
	if res, err := c.Manager.RestartPlugin(uuid); err == plugin.ErrAlreadyEnabledOrDisabled {
		responseHandler(res, err, ctx)
	} else if err != nil {
		responseHandler(res, nil, ctx)
	}
	responseHandler("plugin restart ok", err, ctx)

}

// GetDisplay get display info for Displayer plugin.
func (c *PluginAPI) GetDisplay(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	instance, err := c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}
	ctx.JSON(200, instance.GetDisplay(location.Get(ctx)))

}

// GetConfig returns Configurer plugin configuration in YAML format.
func (c *PluginAPI) GetConfig(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	instance, err := c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}

	if aborted := supportOrAbort(ctx, instance, compat.Configurer); aborted {
		return
	}

	ctx.Header("content-type", "application/x-yaml")
	ctx.Writer.Write(conf.Config)

}

// UpdateConfig updates Configurer plugin configuration in YAML format.
func (c *PluginAPI) UpdateConfig(ctx *gin.Context) {
	uuid := c.buildUUID(ctx)
	conf, err := c.DB.GetPlugin(uuid)
	if success := successOrAbort(ctx, 500, err); !success {
		return
	}
	if conf == nil {
		ctx.AbortWithError(404, errors.New("unknown plugin"))
		return
	}
	instance, err := c.Manager.Instance(uuid)
	if err != nil {
		ctx.AbortWithError(404, errors.New("plugin instance not found"))
		return
	}

	if aborted := supportOrAbort(ctx, instance, compat.Configurer); aborted {
		return
	}

	newConf := instance.DefaultConfig()
	newConfBytes, err := ioutil.ReadAll(ctx.Request.Body)
	var jsonConf map[string]string
	err = json.Unmarshal(newConfBytes, &jsonConf)
	if err != nil {
		ctx.AbortWithError(400, errors.New("invalid data"))
		return
	}
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	if err := yaml.Unmarshal([]byte(jsonConf["data"]), newConf); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	if err := instance.ValidateAndSetConfig(newConf); err != nil {
		ctx.AbortWithError(400, err)
		return
	}
	config, _ := yaml.Marshal(instance.GetConfig())
	conf.Config = config
	successOrAbort(ctx, 500, c.DB.UpdatePluginConf(conf))
}

func supportOrAbort(ctx *gin.Context, instance compat.PluginInstance, module compat.Capability) (aborted bool) {
	if compat.HasSupport(instance, module) {
		return false
	}
	ctx.AbortWithError(400, fmt.Errorf("plugin does not support %s", module))
	return true
}
