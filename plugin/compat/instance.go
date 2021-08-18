package compat

import (
	"github.com/NubeDev/flow-framework/model"
	"github.com/gin-gonic/gin"
	"net/url"
)

// Capability is a capability the plugin provides.
type Capability string

const (
	// Messenger sends notifications.
	Messenger = Capability("messenger")
	// Configurer are consigurables.
	Configurer = Capability("configurer")
	// Storager stores data.
	Storager = Capability("storager")
	// Webhooker registers webhooks.
	Webhooker = Capability("webhooker")
	// Displayer displays instructions.
	Displayer = Capability("displayer")
)

// PluginInstance is an encapsulation layer of plugin instances of different backends.
type PluginInstance interface {
	GetNetworks() ([]*model.Network, error)
	GetNetwork(id string) error
	Enable() error
	Disable() error

	// GetDisplay see Displayer
	GetDisplay(location *url.URL) string

	// DefaultConfig see Configurer
	DefaultConfig() interface{}
	// ValidateAndSetConfig see Configurer
	ValidateAndSetConfig(c interface{}) error

	// SetMessageHandler see Messenger#SetMessageHandler
	SetMessageHandler(h MessageHandler)

	// RegisterWebhook see Webhooker#RegisterWebhook
	RegisterWebhook(basePath string, mux *gin.RouterGroup)

	// SetStorageHandler see Storager#SetStorageHandler.
	SetStorageHandler(handler StorageHandler)

	//Supports Returns the supported modules, f.ex. storager
	Supports() Capabilities
}

// HasSupport tests a PluginInstance for a capability.
func HasSupport(p PluginInstance, toCheck Capability) bool {
	for _, module := range p.Supports() {
		if module == toCheck {
			return true
		}
	}
	return false
}

// Capabilities is a slice of module.
type Capabilities []Capability

// Strings converts []Module to []string.
func (m Capabilities) Strings() []string {
	var result []string
	for _, module := range m {
		result = append(result, string(module))
	}
	return result
}

// MessageHandler see plugin.MessageHandler.
type MessageHandler interface {
	// SendMessage see plugin.MessageHandler
	SendMessage(msg Message) error
}

// StorageHandler see plugin.StorageHandler.
type StorageHandler interface {
	Save(b []byte) error
	Load() ([]byte, error)
	GetNet() ([]*model.Network, error)
}

// Message describes a message to be sent by MessageHandler#SendMessage.
type Message struct {
	Message  			string
	MessageType   		model.MessageType
	IsProtocol    		bool
	DriverType    		model.DriverType
	ProtocolType 		model.ProtocolType
	Protocol			model.Protocol
	WriteableNetwork 	model.WriteableNetwork
	Title    			string
	Priority 			int
	Extras   			map[string]interface{}
}
