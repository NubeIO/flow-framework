package mqttclient

import (
	log "github.com/sirupsen/logrus"
)

var m *Client

// InternalMQTT internal non-secure mqtt connection
// for plugins use the plugin path as the topic
func InternalMQTT(ip string) (bool, error) {
	c, err := NewClient(ClientOptions{
		Servers: []string{ip},
	})
	if err != nil {
		log.Println(err, "CONNECT to broker")
		return false, err
	}
	log.Println(err, "CONNECT to broker")
	m = c
	err = c.Connect()
	if err != nil {
		return false, err
	}
	return c.IsConnected(), nil
}

func GetMQTT() (*Client, bool) {
	return m, m.IsConnected()
}
