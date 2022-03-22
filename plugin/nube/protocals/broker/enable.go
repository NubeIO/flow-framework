package main

import (
	"github.com/fhmq/hmq/broker"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

// Enable implements plugin.Plugin
func (i *Instance) Enable() error {
	i.setUUID()
	if !i.brokerEnabled {
		go i.enableBroker()
	}
	i.brokerEnabled = true
	i.enabled = true
	return nil
}

// Disable implements plugin.Disable
func (i *Instance) Disable() error {
	i.enabled = false
	return nil
}

func (i *Instance) enableBroker() {
	port := "1883"
	if i.config.Port != "" {
		port = i.config.Port
	}
	HttpPort := "8099"
	if i.config.HttpPort != "" {
		HttpPort = i.config.HttpPort
	}
	os.Args = []string{"-port", port, "-httpport", HttpPort}
	config, err := broker.ConfigureConfig(os.Args)
	if err != nil {
		log.Error("configure broker config error: ", err)
	}
	b, err := broker.NewBroker(config)
	if err != nil {
		log.Error("New Broker error: ", err)
	}
	b.Start()
	s := waitForSignal()
	log.Println("signal received, broker closed.", s)
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
