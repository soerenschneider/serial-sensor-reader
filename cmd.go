package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type sensorReaderConfig struct {
	mqttUri           string
	mqttTopic         string
	sensorName        string
	sensorLocation    string
	serialDevice      string
	baudRate          int
	loglevel          string
	prometheusAddress string
}

func (c *sensorReaderConfig) getTopic() string {
	if strings.Contains(c.mqttTopic, "%s") {
		return fmt.Sprintf(c.mqttTopic, c.sensorLocation)
	}

	return c.mqttTopic
}

func (c *sensorReaderConfig) getMqttClientId() string {
	return fmt.Sprintf("sensorreader-%s-%s", c.sensorName, c.sensorLocation)
}

func (c *sensorReaderConfig) printConfig() {
	log.Infof("Started using configuration:")
	log.Infof("mqttUri=%s", c.mqttUri)
	log.Infof("mqttTopic=%s", c.mqttTopic)
	log.Infof("sensorName=%s", c.sensorName)
	log.Infof("sensorLocation=%s", c.sensorLocation)
	log.Infof("serialDevice=%s", c.serialDevice)
	log.Infof("baudRate=%d", c.baudRate)
	log.Infof("loglevel=%s", c.loglevel)
	log.Infof("prometheusAddress=%s", c.prometheusAddress)
	log.Info("")
}

func parseArgs() *sensorReaderConfig {
	parser := argparse.NewParser("serial sensor data reader", "reads data from a serial sensor and publishes it via mqtt")

	mqttUri := parser.String("m", "mqtt-host", &argparse.Options{Required: true, Help: "uri of the mqtt broker. Example: mqtt://remote-host:1883"})
	mqttTopic := parser.String("t", "mqtt-topic", &argparse.Options{Default: "sensors/light/%s", Help: "MQTT topic to send the data to"})
	serialDevice := parser.String("s", "serial-device", &argparse.Options{Default: "/dev/ttyUSB0", Required: false, Help: "The serial device to read sensor data from"})
	baudRate := parser.Int("b", "baud-rate", &argparse.Options{Default: 9600, Required: false, Help: "The baud rate to use for the serial communication"})
	loglevel := parser.Selector("", "loglevel", []string{"INFO", "DEBUG", "WARN"}, &argparse.Options{Default: "INFO", Help: "Debugging loglevel to use"})
	prometheusAddress := parser.String("", "prometheus-address", &argparse.Options{Default: ":9191", Required: false, Help: "Address to use"})
	sensorName := parser.String("n", "sensor-name", &argparse.Options{Required: true, Help: "A descriptive name of the sensor that is read"})
	sensorLocation := parser.String("l", "sensor-location", &argparse.Options{Required: true, Help: "The location of the sensor"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return &sensorReaderConfig{
		mqttUri:           *mqttUri,
		mqttTopic:         *mqttTopic,
		serialDevice:      *serialDevice,
		baudRate:          *baudRate,
		loglevel:          *loglevel,
		prometheusAddress: *prometheusAddress,
		sensorName:        *sensorName,
		sensorLocation:    *sensorLocation,
	}
}
