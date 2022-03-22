package compat

import (
	"github.com/NubeIO/flow-framework/plugin/plugin-api"

	//papiv1 "github.com/NubeIO/flow-framework/plugin/plugin-api"

	"net/url"
	//papiv1 "github.com/NubeIO/flow-framework/plugin/plugin-api"
	"github.com/gin-gonic/gin"
)

// PluginV1 is an abstraction of a plugin written in the v1 plugin API. Exported for testing purposes only.
type PluginV1 struct {
	Info        plugin.Info
	Constructor func() plugin.Plugin
}

// APIVersion returns the API version.
func (c PluginV1) APIVersion() string {
	return "v1"
}

// PluginInfo implements compat/Plugin.
func (c PluginV1) PluginInfo() Info {
	return Info{
		Version:     c.Info.Version,
		Author:      c.Info.Author,
		Name:        c.Info.Name,
		Website:     c.Info.Website,
		Description: c.Info.Description,
		License:     c.Info.License,
		ModulePath:  c.Info.ModulePath,
		HasNetwork:  c.Info.HasNetwork,
	}
}

// NewPluginInstance implements compat/Plugin.
func (c PluginV1) NewPluginInstance() PluginInstance {
	instance := c.Constructor()

	compat := &PluginV1Instance{
		instance: instance,
	}

	if displayer, ok := instance.(plugin.Displayer); ok {
		compat.displayer = displayer
	}

	if messenger, ok := instance.(plugin.Messenger); ok {
		compat.messenger = messenger
	}

	if configurer, ok := instance.(plugin.Configurer); ok {
		compat.configurer = configurer
	}

	if storager, ok := instance.(plugin.Storager); ok {
		compat.storager = storager
	}

	if webhooker, ok := instance.(plugin.Webhooker); ok {
		compat.webhooker = webhooker
	}

	return compat
}

// PluginV1Instance is an adapter for plugin using v1 API.
type PluginV1Instance struct {
	instance   plugin.Plugin
	messenger  plugin.Messenger
	configurer plugin.Configurer
	storager   plugin.Storager
	webhooker  plugin.Webhooker
	displayer  plugin.Displayer
}

// DefaultConfig see papiv1.Configurer.
func (c *PluginV1Instance) DefaultConfig() interface{} {
	if c.configurer != nil {
		return c.configurer.DefaultConfig()
	}
	return struct{}{}
}

func (c *PluginV1Instance) GetConfig() interface{} {
	if c.configurer != nil {
		return c.configurer.GetConfig()
	}
	return struct{}{}
}

// ValidateAndSetConfig see papiv1.Configurer.
func (c *PluginV1Instance) ValidateAndSetConfig(config interface{}) error {
	if c.configurer != nil {
		return c.configurer.ValidateAndSetConfig(config)
	}
	return nil
}

// GetDisplay see papiv1.Displayer.
func (c *PluginV1Instance) GetDisplay(location *url.URL) plugin.Response {
	if c.displayer != nil {
		return c.displayer.GetDisplay(location)
	}
	m := plugin.Help{
		Name:      "name",
		IsNetwork: false,
	}
	r := plugin.Response{
		Details: m,
	}
	return r
}

// SetMessageHandler see papiv1.Messenger.
func (c *PluginV1Instance) SetMessageHandler(h MessageHandler) {
	if c.messenger != nil {
		c.messenger.SetMessageHandler(&PluginV1MessageHandler{WrapperHandler: h})
	}
}

// RegisterWebhook see papiv1.Webhooker.
func (c *PluginV1Instance) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	if c.webhooker != nil {
		c.webhooker.RegisterWebhook(basePath, mux)
	}
}

// SetStorageHandler see papiv1.Storager.
func (c *PluginV1Instance) SetStorageHandler(handler StorageHandler) {
	if c.storager != nil {
		c.storager.SetStorageHandler(&PluginV1StorageHandler{WrapperHandler: handler})
	}
}

// Supports returns a slice of capabilities the plugin instance provides.
func (c *PluginV1Instance) Supports() Capabilities {
	modules := Capabilities{}
	if c.configurer != nil {
		modules = append(modules, Configurer)
	}
	if c.displayer != nil {
		modules = append(modules, Displayer)
	}
	if c.messenger != nil {
		modules = append(modules, Messenger)
	}
	if c.storager != nil {
		modules = append(modules, Storager)
	}
	if c.webhooker != nil {
		modules = append(modules, Webhooker)
	}
	return modules
}

// PluginV1MessageHandler is an adapter for messenger plugin handler using v1 API.
type PluginV1MessageHandler struct {
	WrapperHandler MessageHandler
}

// SendMessage implements papiv1.MessageHandler.
func (c *PluginV1MessageHandler) SendMessage(msg plugin.Message) error {
	return c.WrapperHandler.SendMessage(Message{
		Message:  msg.Message,
		Priority: msg.Priority,
		Title:    msg.Title,
		Extras:   msg.Extras,
	})
}

// Enable implements wrapper.Plugin.
func (c *PluginV1Instance) Enable() error {
	return c.instance.Enable()
}

// Disable implements wrapper.Plugin.
func (c *PluginV1Instance) Disable() error {
	return c.instance.Disable()
}

// PluginV1StorageHandler is a wrapper for v1 storage handler.
type PluginV1StorageHandler struct {
	WrapperHandler StorageHandler
}

// Save implements wrapper.Storager.
func (c *PluginV1StorageHandler) Save(b []byte) error {
	return c.WrapperHandler.Save(b)
}

// Load implements wrapper.Storager.
func (c *PluginV1StorageHandler) Load() ([]byte, error) {
	return c.WrapperHandler.Load()
}
