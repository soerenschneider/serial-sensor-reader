package main

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

const (
	quiesceMilliseconds   = 5000
	connectTimeoutSec     = 180
)

type mqttBackend struct {
	uri      *url.URL
	client   mqtt.Client
	clientId string
}

func NewMqttBackend(conf *sensorReaderConfig) (*mqttBackend, error) {
	if nil == conf {
		return nil, errors.New("nil sensorReaderConfig supplied")
	}

	mqttBackend := &mqttBackend{
		uri:    conf.mqttUri,
	}

	err := mqttBackend.connect(conf)
	return mqttBackend, err
}

func (b *mqttBackend) connect(conf *sensorReaderConfig) error {
	opts := createClientOptions(conf, b.uri)

	client := mqtt.NewClient(opts)
	log.Info("Trying to connect to broker")
	token := client.Connect()
	for !token.WaitTimeout(connectTimeoutSec * time.Second) {
	}
	if err := token.Error(); err != nil {
		return err
	}

	b.client = client
	return nil
}

func createClientOptions(conf *sensorReaderConfig, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf(uri.Host))

	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)

	opts.SetClientID(conf.getMqttClientId())

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Infof("(Re-)connected to MQTT broker %s...", uri.Host)
		metricMqttConnectionsEstablished.WithLabelValues(conf.sensorName, conf.sensorLocation).Inc()
	})

	opts.SetConnectionLostHandler(func(mqtt.Client, error) {
		log.Error("Lost connection to broker")
		metricMqttConnectionsLost.WithLabelValues(conf.sensorName, conf.sensorLocation).Inc()
	})

	return opts
}

func (b *mqttBackend) Disconnect() {
	b.client.Disconnect(quiesceMilliseconds)
}

func (b *mqttBackend) Publish(topic, message string) error {
	log.Debugf("Publishing message to topic %s", topic)
	token := b.client.Publish(topic, 0, false, message)
	return token.Error()
}
