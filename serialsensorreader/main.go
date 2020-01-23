package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type serialSensorReader struct {
	mqtt   *mqttBackend
	sensor *serialReader
	config *sensorReaderConfig
}

func main() {
	config := parseArgs()

	setupLogging(config)
	config.printConfig()
	initializeMetrics()
	go startPrometheusMetricsServer(config.prometheusAddress)

	b := serialSensorReader{config: config}
	b.mqtt = createMqttBackend(config)
	b.sensor = createSerialReader(config)

	output := make(chan string)
	go b.sensor.ReadSensorData(output)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-done
		log.Info("Received SIGTERM, quitting gracefully")
		b.cleanup(output)
	}()
	
	topic := config.getTopic()
	for dataPoint := range output {
		b.handleSensorDatapoint(dataPoint, topic)
	}
}

func (b *serialSensorReader) cleanup(output chan string) {
	b.sensor.Interrupt()
	log.Info("Stopped reading from serial device")
	b.mqtt.Disconnect()
	log.Info("Disconnected from mqtt")
	close(output)
	log.Info("Closed channel, bye")
}

func (b *serialSensorReader) handleSensorDatapoint(dataPoint string, topic string) error {
	err := b.mqtt.Publish(topic, dataPoint)
	if err != nil {
		log.Errorf("Could not publish datapoint: %s", err.Error())
	}

	numericValue, err := strconv.ParseFloat(dataPoint, 64)
	if err != nil {
		log.Warningf("Received invalid value %s", dataPoint)
		return err
	}

	metricSensorData.WithLabelValues(b.config.sensorName, b.config.sensorLocation).Set(numericValue)
	metricSensorDataCnt.WithLabelValues(b.config.sensorName, b.config.sensorLocation).Inc()
	metricSensorLastMeasurement.WithLabelValues(b.config.sensorName, b.config.sensorLocation).SetToCurrentTime()

	return nil
}

func setupLogging(conf *sensorReaderConfig) {
	switch conf.loglevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func createMqttBackend(config *sensorReaderConfig) *mqttBackend {
	log.Info("Creating mqtt backend...")
	mqtt, err := NewMqttBackend(config)
	if err != nil {
		log.Fatalf("Couldn't setup mqtt backend: %s", err.Error())
	}
	log.Info("Mqtt backend initialized")
	return mqtt
}

func createSerialReader(config *sensorReaderConfig) *serialReader {
	log.Infof("Starting to read from serial device %s with baudrate %d", config.serialDevice, config.baudRate)
	r, err := NewSerialReader(config.serialDevice, config.baudRate)
	if err != nil {
		log.Fatal(err)
	}
	return r
}