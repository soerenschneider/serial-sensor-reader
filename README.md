# serial-sensor-reader

Reads sensor data over a serial connection, e.g. an Arduino and publishes the 
read sensor data via MQTT. Also, vital metrics are exposed via Prometheus server.

## installation

To build it on your local machine just check out the repository and invoke
```
make build
``` 

To cross-compile for Raspberry PIs devices, just use the raspberry target:
```
make raspberry
```

## usage

```
usage: serial sensor data reader [-h|--help] -m|--mqtt-host "<value>"
                                 [-t|--mqtt-topic "<value>"]
                                 [-s|--serial-device "<value>"] [-b|--baud-rate
                                 <integer>] [--loglevel (INFO|DEBUG|WARN)]
                                 [--prometheus-address "<value>"]
                                 -n|--sensor-name "<value>"
                                 -l|--sensor-location "<value>"

                                 reads data from a serial sensor and publishes
                                 it via mqtt

Arguments:

  -h  --help                Print help information
  -m  --mqtt-host           uri of the mqtt broker. Example:
                            mqtt://remote-host:1883
  -t  --mqtt-topic          MQTT topic to send the data to. Default:
                            sensors/light/%s
  -s  --serial-device       The serial device to read sensor data from.
                            Default: /dev/ttyUSB0
  -b  --baud-rate           The baud rate to use for the serial communication.
                            Default: 9600
      --loglevel            Debugging loglevel to use. Default: INFO
      --prometheus-address  Address to use. Default: :9191
  -n  --sensor-name         A descriptive name of the sensor that is read
  -l  --sensor-location     The location of the sensor
```

## example

```
 ./build/serial-sensor-reader -m mqtt://localhost:1883 -t sensors/photo/wohnzimmer -n photoresistor -l wohnzimmer --loglevel DEBUG
INFO[2020-01-22T19:02:41+01:00] Started using configuration:                 
INFO[2020-01-22T19:02:41+01:00] mqttUri=mqtt://localhost:1883                
INFO[2020-01-22T19:02:41+01:00] mqttTopic=sensors/photoresistor/wohnzimmer           
INFO[2020-01-22T19:02:41+01:00] sensorName=light                             
INFO[2020-01-22T19:02:41+01:00] sensorLocation=wohnzimmer                    
INFO[2020-01-22T19:02:41+01:00] serialDevice=/dev/ttyUSB0                    
INFO[2020-01-22T19:02:41+01:00] baudRate=9600                                
INFO[2020-01-22T19:02:41+01:00] loglevel=DEBUG                               
INFO[2020-01-22T19:02:41+01:00] prometheusAddress=:9191                      
INFO[2020-01-22T19:02:41+01:00]                                              
INFO[2020-01-22T19:02:41+01:00] Creating mqtt backend...                     
INFO[2020-01-22T19:02:41+01:00] Mqtt backend initialized                     
INFO[2020-01-22T19:02:41+01:00] Starting to read from serial device /dev/ttyUSB0 with baudrate 9600 
INFO[2020-01-22T19:02:41+01:00] (Re-)connected to MQTT broker localhost:1883... 
DEBU[2020-01-22T19:02:42+01:00] Read value 901                               
DEBU[2020-01-22T19:02:42+01:00] Publishing message to topic sensors/photo/wohnzimmer 
DEBU[2020-01-22T19:02:43+01:00] Read value 901                               
DEBU[2020-01-22T19:02:43+01:00] Publishing message to topic sensors/photo/wohnzimmer 
^CINFO[2020-01-22T19:02:44+01:00] Received SIGTERM, quitting gracefully        
INFO[2020-01-22T19:02:44+01:00] Stopped reading from serial device           
INFO[2020-01-22T19:02:44+01:00] Disconnected from mqtt                       
INFO[2020-01-22T19:02:44+01:00] Closed channel, bye 
```
