package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	metricMqttMessagesSent           *prometheus.CounterVec
	metricMqttConnectionsEstablished *prometheus.CounterVec
	metricMqttConnectionsLost        *prometheus.CounterVec
	metricSensorData                 *prometheus.GaugeVec
	metricSensorDataCnt				 *prometheus.CounterVec
	metricSensorLastMeasurement      *prometheus.GaugeVec
)

func initializeMetrics() {
	metricMqttMessagesSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "serialsensordata_mqtt_messages_sent_total",
		Help: "The total number of message sent via MQTT",
	}, []string {"sensor", "location"})

	metricMqttConnectionsEstablished = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "serialsensordata_mqtt_connections_established_total",
		Help: "The total number of MQTT connections established",
	}, []string {"sensor", "location"})

	metricMqttConnectionsLost = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "serialsensordata_mqtt_connections_lost_total",
		Help: "The total number of MQTT connections lost",
	}, []string {"sensor", "location"})

	metricSensorData = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "serialsensordata_sensor_data",
		Help: "Current sensor value",
	}, []string {"sensor", "location"})

	metricSensorLastMeasurement = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "serialsensordata_sensor_last_measurement_timestamp",
		Help: "Timestamp of the last sensor read operation",
	}, []string {"sensor", "location"})

	metricSensorDataCnt = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "serialsensordata_sensor_data_count_total",
		Help: "The total number of read sensor data operations",
	}, []string {"sensor", "location"})
}

func startPrometheusMetricsServer(address string) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("couldn't start prometheus listener: %s", err.Error())
	}
}