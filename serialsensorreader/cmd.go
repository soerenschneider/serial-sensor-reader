package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"strings"
	"text/template"
)

const (
	defaultMqttTopic = "sensors/{{ .Name }}/{{ .Location }}"
	defaultSerialDevice = "/dev/ttyUSB0"
	defaultBaudRate = 9600
	defaultPrometheusAddress = ":9191"
	defaultLoglevel = "INFO"
)

type sensorReaderConfig struct {
	mqttUri           *url.URL
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

func (c *sensorReaderConfig) printParsedValues() {
	log.Infof("Started using configuration:")
	log.Info("---------")
	log.Infof("mqttUri=%s", c.mqttUri)
	log.Infof("mqttTopic=%s", c.mqttTopic)
	log.Infof("sensorName=%s", c.sensorName)
	log.Infof("sensorLocation=%s", c.sensorLocation)
	log.Infof("serialDevice=%s", c.serialDevice)
	log.Infof("baudRate=%d", c.baudRate)
	log.Infof("loglevel=%s", c.loglevel)
	log.Infof("prometheusAddress=%s", c.prometheusAddress)
	log.Info("---------")
}

func parseArgs() *sensorReaderConfig {
	parser := argparse.NewParser("serial sensor data reader", "reads data from a serial sensor and publishes it via mqtt")

	mqttUri := parser.String("m", "mqtt-host", &argparse.Options{Required: true, Help: "uri of the mqtt broker. Example: mqtt://remote-host:1883"})
	mqttTopic := parser.String("t", "mqtt-topic", &argparse.Options{Default: defaultMqttTopic, Help: "MQTT topic to send the data to"})
	serialDevice := parser.String("s", "serial-device", &argparse.Options{Default: defaultSerialDevice, Required: false, Help: "The serial device to read sensor data from"})
	baudRate := parser.Int("b", "baud-rate", &argparse.Options{Default: defaultBaudRate, Required: false, Help: "The baud rate to use for the serial communication"})
	loglevel := parser.Selector("", "loglevel", []string{"INFO", "DEBUG", "WARN"}, &argparse.Options{Default: defaultLoglevel, Help: "Debugging loglevel to use"})
	prometheusAddress := parser.String("", "prometheus-address", &argparse.Options{Default: defaultPrometheusAddress, Required: false, Help: "Address to use"})
	sensorName := parser.String("n", "sensor-name", &argparse.Options{Required: true, Help: "A descriptive name of the sensor that is read"})
	sensorLocation := parser.String("l", "sensor-location", &argparse.Options{Required: true, Help: "The location of the sensor"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	
	topic := getTopic(*mqttTopic, *sensorName, *sensorLocation)
	fmt.Println(topic)
	uri, err := parseUrl(*mqttUri)
	if err != nil {
		log.Fatalf("Invalid broker url supplied: %s", err.Error())
	}

	return &sensorReaderConfig{
		mqttUri:           uri,
		mqttTopic:         topic,
		serialDevice:      *serialDevice,
		baudRate:          *baudRate,
		loglevel:          *loglevel,
		prometheusAddress: *prometheusAddress,
		sensorName:        *sensorName,
		sensorLocation:    *sensorLocation,
	}
}

func getTopic(templateString, sensorName, sensorLocation string) string {
	tmpl, err := template.New("mqttTopic").Parse(templateString)
	buf := &bytes.Buffer{}

	data := struct {
		Name string
		Location string
	}{
		sensorName,
		sensorLocation,
	}
	
	err = tmpl.Execute(buf, data)
	if err != nil {
		fmt.Printf("Couldn't process template: %s", err.Error())
		os.Exit(1)
	}
	return buf.String()
}

func parseUrl(rawUrl string) (*url.URL, error) {
	uri, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	if len(uri.Port()) == 0 {
		return nil, errors.New("can not connect to broker, no port supplied")
	}

	if len(uri.Scheme) == 0 {
		return nil, errors.New("can not connect to broker, no scheme supplied")
	}

	return uri, nil
}